package message_repository

import (
	message_model "github.com/evolution-foundation/evolution-go/pkg/message/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MessageRepository interface {
	InsertMessage(message message_model.Message) error
	GetMessageByID(messageID string) (*message_model.Message, error)
	FindRecentMessages(filter MessageHistoryFilter) ([]message_model.Message, error)
	DeleteAllMessages() (int64, error)
	GetLatestMessageID(source string) (string, string, error)
}

type MessageHistoryFilter struct {
	InstanceID string
	MessageID  string
	Chat       string
	MediaType  string
	Query      string
	Limit      int
	IncludeRaw bool
}

type messageRepository struct {
	db *gorm.DB
}

func messageUpdateColumns(message message_model.Message) []string {
	updates := []string{"timestamp", "status", "source"}
	if hasRichPayload(message) {
		updates = append(updates,
			"instance_id",
			"chat_jid",
			"sender_jid",
			"participant_jid",
			"chat_name",
			"sender_name",
			"from_me",
			"is_group",
			"text",
			"media_type",
			"media_url",
			"mimetype",
			"raw_json",
		)
	}
	if len(message.Referral) > 0 {
		updates = append(updates, "referral")
	}

	return updates
}

func hasRichPayload(message message_model.Message) bool {
	return message.InstanceID != "" ||
		message.ChatJID != "" ||
		message.SenderJID != "" ||
		message.ParticipantJID != "" ||
		message.ChatName != "" ||
		message.SenderName != "" ||
		message.Text != "" ||
		message.MediaType != "" ||
		message.MediaURL != "" ||
		message.Mimetype != "" ||
		len(message.RawJSON) > 0
}

func (m *messageRepository) InsertMessage(message message_model.Message) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}},
		DoUpdates: clause.AssignmentColumns(messageUpdateColumns(message)),
	}).Create(&message).Error
}

func (m *messageRepository) GetMessageByID(messageID string) (*message_model.Message, error) {
	var message message_model.Message
	err := m.db.Where("message_id = ?", messageID).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &message, nil
}

func (m *messageRepository) FindRecentMessages(filter MessageHistoryFilter) ([]message_model.Message, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	query := m.db.Model(&message_model.Message{})
	if filter.InstanceID != "" {
		query = query.Where("(instance_id = ? OR instance_id = '' OR instance_id IS NULL)", filter.InstanceID)
	}
	if filter.MessageID != "" {
		query = query.Where("message_id = ?", filter.MessageID)
	}
	if filter.Chat != "" {
		query = query.Where(
			"(chat_jid = ? OR sender_jid = ? OR participant_jid = ? OR source = ?)",
			filter.Chat,
			filter.Chat,
			filter.Chat,
			filter.Chat,
		)
	}
	if filter.MediaType != "" {
		query = query.Where("media_type = ?", filter.MediaType)
	}
	if filter.Query != "" {
		like := "%" + filter.Query + "%"
		query = query.Where(
			"(text ILIKE ? OR chat_name ILIKE ? OR sender_name ILIKE ? OR chat_jid ILIKE ? OR sender_jid ILIKE ? OR participant_jid ILIKE ? OR media_url ILIKE ?)",
			like,
			like,
			like,
			like,
			like,
			like,
			like,
		)
	}
	if !filter.IncludeRaw {
		query = query.Omit("raw_json")
	}

	var messages []message_model.Message
	err := query.Order("timestamp DESC").Order("id DESC").Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *messageRepository) DeleteAllMessages() (int64, error) {
	result := m.db.Exec("DELETE FROM messages")
	return result.RowsAffected, result.Error
}

func (m *messageRepository) GetLatestMessageID(source string) (string, string, error) {
	var message message_model.Message
	err := m.db.Where("source = ?", source).Order("timestamp DESC").First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", nil
		}
		return "", "", err
	}

	return message.MessageID, message.Timestamp, nil
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}
