# SecretBay VPN Configuration Tool

A secure and automated tool for configuring OpenVPN and IKEv2/IPsec VPN servers.

## Features

- Automated configuration of OpenVPN and IKEv2/IPsec VPN servers
- Support for both password and SSH key authentication
- Secure password generation and root password management
- Clean and intuitive web interface
- iOS-compatible VPN configuration (.mobileconfig)
- Automated cleanup of sensitive data
- Docker-based deployment for easy installation

## Prerequisites

- Linux-based system
- Docker and Docker Compose (automatically installed by deploy script)
- Internet connection

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/SecretBayLite.git
cd SecretBayLite
```

2. Run the deployment script:
```bash
./deploy.sh
```

The script will:
- Install Docker and Docker Compose if not present
- Build and start the containers
- Set up the necessary directories and permissions

3. Access the web interface at `http://localhost:3000`

## Usage

1. Open the web interface in your browser
2. Enter the VPN server details:
   - Server IP address
   - Username (root or sudo user)
   - Authentication method (password or SSH key)
   - VPN type (OpenVPN or IKEv2/IPsec)
3. Submit the configuration
4. Download and save the generated configuration files
5. Note the new root password provided for future access

## Security Features

- Automatic root password rotation
- Secure configuration standards for both OpenVPN and IKEv2
- Fail2ban integration for brute-force protection
- Automated cleanup of sensitive logs and temporary files
- JWT-based API authentication

## Architecture

The application consists of two main components:

1. Backend (Go):
   - RESTful API for VPN configuration
   - SSH-based server management
   - Secure credential handling
   - Automated package installation and setup

2. Frontend (Flutter Web):
   - Responsive web interface
   - Real-time validation
   - Secure credential input
   - Configuration file download

## Development

### Backend Development

```bash
cd backend
go mod tidy
go test ./...
go run cmd/server/main.go
```

### Frontend Development

```bash
cd frontend
flutter pub get
flutter run -d chrome
```

## Testing

Run the test suites:

```bash
# Backend tests
cd backend
go test -v ./...

# Frontend tests
cd frontend
flutter test
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Security Considerations

- Always change the JWT secret in production
- Use strong passwords for authentication
- Keep your server's SSH key secure
- Regularly update the system and dependencies
- Monitor server logs for suspicious activity