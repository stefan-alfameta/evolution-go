package user_service

import "testing"

func TestParseUserJIDCanonicalizesPhoneNumbers(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "plain brazilian number",
			in:   "556696701002",
			want: "556696701002@s.whatsapp.net",
		},
		{
			name: "formatted brazilian number",
			in:   "+55 (66) 9670-1002",
			want: "556696701002@s.whatsapp.net",
		},
		{
			name: "already a whatsapp jid",
			in:   "556696701002@s.whatsapp.net",
			want: "556696701002@s.whatsapp.net",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseUserJID(tt.in)
			if !ok {
				t.Fatalf("parseUserJID(%q) failed", tt.in)
			}
			if got.String() != tt.want {
				t.Fatalf("parseUserJID(%q) = %q, want %q", tt.in, got.String(), tt.want)
			}
			if got.User[0] == '+' {
				t.Fatalf("parseUserJID(%q) returned non-canonical user %q", tt.in, got.User)
			}
		})
	}
}
