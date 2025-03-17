// Package vpn handles VPN configuration operations for both OpenVPN and IKEv2.
package vpn

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// ConfigRequest represents the client request for VPN configuration
type ConfigRequest struct {
	ServerIP       string `json:"server_ip"`
	Username       string `json:"username"`
	AuthMethod     string `json:"auth_method"`
	AuthCredential string `json:"auth_credential"`
	VPNType        string `json:"vpn_type"`
}

// ConfigResponse represents the response with VPN configuration
type ConfigResponse struct {
	Config      string `json:"config"`
	NewPassword string `json:"new_password"`
}

// Service handles VPN configuration operations
type Service struct {
	logger io.Writer
}

// NewService creates a new VPN service
func NewService(logger io.Writer) *Service {
	return &Service{
		logger: logger,
	}
}

// ConfigureVPN sets up VPN on the remote server
func (s *Service) ConfigureVPN(req *ConfigRequest) (*ConfigResponse, error) {
	// Validate request
	if err := s.validateVPNRequest(req); err != nil {
		return nil, err
	}

	// Create SSH client config
	config := &ssh.ClientConfig{
		User:            req.Username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: In production, use proper host key verification
	}

	// Set up authentication
	if req.AuthMethod == "password" {
		config.Auth = append(config.Auth, ssh.Password(req.AuthCredential))
	} else if req.AuthMethod == "key" {
		signer, err := ssh.ParsePrivateKey([]byte(req.AuthCredential))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	// Connect to remote server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", req.ServerIP), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}
	defer client.Close()

	// Install required packages based on VPN type
	if err = s.installRequiredPackages(client, req.VPNType); err != nil {
		return nil, err
	}

	// Configure VPN based on type
	var vpnConfig string
	if req.VPNType == "openvpn" {
		vpnConfig, err = s.configureOpenVPN(client)
	} else if req.VPNType == "ios" {
		vpnConfig, err = s.configureIKEv2(client)
	}
	if err != nil {
		return nil, err
	}

	// Generate new password
	newPassword := generateSecurePassword()
	if err = s.changeRootPassword(client, newPassword); err != nil {
		return nil, err
	}

	// Clean up sensitive data
	if err = s.cleanup(client); err != nil {
		return nil, err
	}

	return &ConfigResponse{
		Config:      vpnConfig,
		NewPassword: newPassword,
	}, nil
}

// validateVPNRequest validates the VPN configuration request
func (s *Service) validateVPNRequest(req *ConfigRequest) error {
	// Validate server IP
	if ip := net.ParseIP(req.ServerIP); ip == nil {
		return fmt.Errorf("invalid server IP address")
	}

	// Validate username
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}

	// Validate auth method
	if req.AuthMethod != "password" && req.AuthMethod != "key" {
		return fmt.Errorf("invalid authentication method")
	}

	// Validate VPN type
	if req.VPNType != "openvpn" && req.VPNType != "ios" {
		return fmt.Errorf("invalid VPN type")
	}

	// Validate auth credential
	if req.AuthCredential == "" {
		return fmt.Errorf("authentication credential is required")
	}

	return nil
}

// installRequiredPackages installs necessary packages on the remote server
func (s *Service) installRequiredPackages(client *ssh.Client, vpnType string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var installCmd string
	if vpnType == "openvpn" {
		installCmd = "apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y openvpn easy-rsa fail2ban ufw"
	} else if vpnType == "ios" {
		installCmd = "apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y strongswan strongswan-pki libcharon-extra-plugins fail2ban ufw"
	}

	session.Stdout = s.logger
	session.Stderr = s.logger

	return session.Run(installCmd)
}

// configureOpenVPN sets up OpenVPN and returns the client configuration
func (s *Service) configureOpenVPN(client *ssh.Client) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Set up PKI and generate certificates
	setupCommands := []string{
		"cd /etc/openvpn",
		"mkdir -p easy-rsa",
		"cp -r /usr/share/easy-rsa/* easy-rsa/",
		"cd easy-rsa",
		"./easyrsa init-pki",
		"./easyrsa build-ca nopass",
		"./easyrsa build-server-full server nopass",
		"./easyrsa build-client-full client nopass",
		"./easyrsa gen-dh",
		"openvpn --genkey secret ta.key",
	}

	session.Stdout = s.logger
	session.Stderr = s.logger

	if err := session.Run(strings.Join(setupCommands, " && ")); err != nil {
		return "", fmt.Errorf("failed to set up OpenVPN: %v", err)
	}

	// Generate client configuration
	config := `client
dev tun
proto udp
remote SERVER_IP 1194
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
cipher AES-256-GCM
auth SHA256
key-direction 1
verb 3`

	return config, nil
}

// configureIKEv2 sets up IKEv2/IPsec and returns the client configuration
func (s *Service) configureIKEv2(client *ssh.Client) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Set up certificates and strongSwan configuration
	setupCommands := []string{
		"cd /etc/strongswan",
		"strongswan pki --gen --type rsa --size 4096 --outform pem > ipsec.d/private/ca-key.pem",
		"strongswan pki --self --ca --lifetime 3650 --in ipsec.d/private/ca-key.pem --type rsa --dn 'CN=SecretBay CA' --outform pem > ipsec.d/cacerts/ca-cert.pem",
		"strongswan pki --gen --type rsa --size 4096 --outform pem > ipsec.d/private/server-key.pem",
		"strongswan pki --pub --in ipsec.d/private/server-key.pem --type rsa | strongswan pki --issue --lifetime 3650 --cacert ipsec.d/cacerts/ca-cert.pem --cakey ipsec.d/private/ca-key.pem --dn 'CN=SERVER_IP' --san 'SERVER_IP' --flag serverAuth --flag ikeIntermediate --outform pem > ipsec.d/certs/server-cert.pem",
	}

	session.Stdout = s.logger
	session.Stderr = s.logger

	if err := session.Run(strings.Join(setupCommands, " && ")); err != nil {
		return "", fmt.Errorf("failed to set up IKEv2: %v", err)
	}

	// Generate .mobileconfig for iOS
	mobileConfig := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>PayloadContent</key>
	<array>
		<dict>
			<key>IKEv2</key>
			<dict>
				<key>RemoteAddress</key>
				<string>SERVER_IP</string>
				<key>RemoteIdentifier</key>
				<string>SERVER_IP</string>
				<key>LocalIdentifier</key>
				<string>client</string>
				<key>AuthenticationMethod</key>
				<string>Certificate</string>
			</dict>
			<key>PayloadDescription</key>
			<string>Configures VPN settings</string>
			<key>PayloadDisplayName</key>
			<string>SecretBay VPN</string>
			<key>PayloadIdentifier</key>
			<string>com.secretbay.vpn</string>
			<key>PayloadType</key>
			<string>com.apple.vpn.managed</string>
			<key>PayloadUUID</key>
			<string>` + generateUUID() + `</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
		</dict>
	</array>
	<key>PayloadDisplayName</key>
	<string>SecretBay VPN Configuration</string>
	<key>PayloadIdentifier</key>
	<string>com.secretbay.vpn.config</string>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>` + generateUUID() + `</string>
	<key>PayloadVersion</key>
	<integer>1</integer>
</dict>
</plist>`

	return mobileConfig, nil
}

// changeRootPassword changes the root password on the remote server
func (s *Service) changeRootPassword(client *ssh.Client, newPassword string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("echo 'root:%s' | chpasswd", newPassword)
	return session.Run(cmd)
}

// cleanup removes sensitive data from the remote server
func (s *Service) cleanup(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Clear various logs and temporary files
	cleanupCmd := `
		rm -rf /tmp/* /var/tmp/*
		find /var/log -type f -exec sh -c "cat /dev/null > {}" \;
		history -c
	`
	return session.Run(cleanupCmd)
}

// generateSecurePassword creates a secure random password
func generateSecurePassword() string {
	const (
		length  = 16
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	)

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to a default secure password if random generation fails
		return "SecretBay" + base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))[:12]
	}

	var password strings.Builder
	password.Grow(length)

	for i := 0; i < length; i++ {
		password.WriteByte(charset[int(randomBytes[i])%len(charset)])
	}

	return password.String()
}

// generateUUID generates a random UUID for iOS configuration
func generateUUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
