package vpn

import (
	"bytes"
	"testing"
)

func TestValidateVPNRequest(t *testing.T) {
	logger := &bytes.Buffer{}
	service := NewService(logger)

	tests := []struct {
		name    string
		req     *ConfigRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &ConfigRequest{
				ServerIP:       "192.168.1.1",
				Username:       "root",
				AuthMethod:     "password",
				AuthCredential: "password123",
				VPNType:        "openvpn",
			},
			wantErr: false,
		},
		{
			name: "invalid IP",
			req: &ConfigRequest{
				ServerIP:       "invalid-ip",
				Username:       "root",
				AuthMethod:     "password",
				AuthCredential: "password123",
				VPNType:        "openvpn",
			},
			wantErr: true,
		},
		{
			name: "invalid auth method",
			req: &ConfigRequest{
				ServerIP:       "192.168.1.1",
				Username:       "root",
				AuthMethod:     "invalid",
				AuthCredential: "password123",
				VPNType:        "openvpn",
			},
			wantErr: true,
		},
		{
			name: "invalid VPN type",
			req: &ConfigRequest{
				ServerIP:       "192.168.1.1",
				Username:       "root",
				AuthMethod:     "password",
				AuthCredential: "password123",
				VPNType:        "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateVPNRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateVPNRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateSecurePassword(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "generates valid password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password := generateSecurePassword()
			if len(password) < 16 {
				t.Errorf("generateSecurePassword() generated password length = %d, want at least 16", len(password))
			}

			// Check if password contains at least one character from each required set
			hasUpper := false
			hasLower := false
			hasNumber := false
			hasSpecial := false

			for _, c := range password {
				switch {
				case c >= 'A' && c <= 'Z':
					hasUpper = true
				case c >= 'a' && c <= 'z':
					hasLower = true
				case c >= '0' && c <= '9':
					hasNumber = true
				case c == '!' || c == '@' || c == '#' || c == '$' || c == '%' || c == '^' || c == '&' || c == '*' || c == '(' || c == ')' || c == '-' || c == '_' || c == '=' || c == '+':
					hasSpecial = true
				}
			}

			if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
				t.Errorf("generateSecurePassword() generated password does not meet complexity requirements")
			}
		})
	}
}

func TestGenerateUUID(t *testing.T) {
	uuid1 := generateUUID()
	uuid2 := generateUUID()

	if uuid1 == uuid2 {
		t.Error("generateUUID() generated identical UUIDs")
	}

	if len(uuid1) != 36 {
		t.Errorf("generateUUID() generated UUID length = %d, want 36", len(uuid1))
	}

	// Check UUID format (8-4-4-4-12)
	parts := bytes.Split([]byte(uuid1), []byte("-"))
	if len(parts) != 5 {
		t.Errorf("generateUUID() generated UUID has %d parts, want 5", len(parts))
	}

	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			t.Errorf("generateUUID() part %d has length %d, want %d", i, len(part), expectedLengths[i])
		}
	}
}
