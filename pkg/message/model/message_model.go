package message_model

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	Id             string          `json:"id" gorm:"type:uuid;primaryKey"`
	MessageID      string          `json:"message_id" gorm:"unique"`
	Timestamp      string          `json:"timestamp"`
	Status         string          `json:"status"`
	Source         string          `json:"source"`
	InstanceID     string          `json:"instance_id,omitempty" gorm:"column:instance_id;index"`
	ChatJID        string          `json:"chat_jid,omitempty" gorm:"column:chat_jid;index"`
	SenderJID      string          `json:"sender_jid,omitempty" gorm:"column:sender_jid;index"`
	ParticipantJID string          `json:"participant_jid,omitempty" gorm:"column:participant_jid;index"`
	ChatName       string          `json:"chat_name,omitempty"`
	SenderName     string          `json:"sender_name,omitempty"`
	FromMe         bool            `json:"from_me"`
	IsGroup        bool            `json:"is_group"`
	Text           string          `json:"text,omitempty"`
	MediaType      string          `json:"media_type,omitempty" gorm:"index"`
	MediaURL       string          `json:"media_url,omitempty"`
	Mimetype       string          `json:"mimetype,omitempty"`
	RawJSON        json.RawMessage `json:"raw_json,omitempty" gorm:"type:jsonb"`
	Referral       json.RawMessage `json:"referral,omitempty" gorm:"type:jsonb"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = uuid.New().String()
	return
}
