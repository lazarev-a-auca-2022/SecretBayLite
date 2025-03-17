Project: Server Application for Automatic Configuration of Remote VPN Servers
Name: SecretBay
Objective: Create a server application that:

    Accepts requests from client applications with data for a remote server.
    Automatically configures VPN on a remote server running Ubuntu.
    Ensures the security of the remote server configuration.
    Returns a ready VPN configuration to the client for OpenVPN and iOS VPN.
    Deletes all client data after completing the configuration (including changing the password).

1. Functional Requirements

1.0.1. It must work on a separate server even when hosting locally. 
1.0.2. An .sh script to deploy the project on any machine.

1.1. Client Request

The server must accept requests with JSON data containing:

    server_ip (string): The IP address of the remote server.
    username (string): The username for login (default is root).
    auth_method (string): The authentication method (password or key).
    auth_credential (string): The password or SSH key.
    vpn_type (string): The type of VPN (iOS VPN or OpenVPN).

1.2. Server Configuration

Upon receiving the data, the server must:

    Connect to the remote server via SSH.
    Install the necessary VPN software:
        For OpenVPN: openvpn, easyrsa.
        For iOS native VPN: StrongSwan.
    Configure the VPN:
        Generate keys and configurations.
        Enable traffic routing.
        Disable log storage.
    Configure security:
        Install fail2ban to protect against attacks.
        Disable unnecessary services and processes.
        Change the root password of the user.

1.3. Return Data to the Client

The server must return the following to the client:

    A VPN configuration file for connection (in .mobileconfig format for iPhone or .ovpn format for OpenVPN) and the new password for the remote server.

1.4. Data Deletion

After completing all tasks, the server must:

    Delete temporary data used for configuration.
    Remove client information (e.g., IP address, SSH keys) from the server.

2. Non-Functional Requirements
2.1. Performance

    The server must handle up to 50 concurrent requests.
    Maximum request processing time: 120 seconds.
    The server must provide the client with progress updates during execution.

2.2. Security

    All communications with the API must be protected via HTTPS.
    Client authentication via tokens (e.g., JWT).
    Task execution must be isolated for different clients.

2.3. Reliability

    In case of an error, execution must be interrupted, and the client must be returned an appropriate error message.
    Error logs must be recorded on the server.

2.4. Scalability

    The application must be deployable in a Docker container for easy scaling.

Programming Language: Go (Golang)

Libraries:

    For HTTP: net/http or gorilla/mux.
    For SSH: golang.org/x/crypto/ssh. VPN services: StrongSwan (IKEv2). Containerization: Docker. Server OS: Ubuntu 22.

Development, Writing, Formatting, and Commenting Guidelines:

The code should follow Google's Go Style Guide:
https://google.github.io/styleguide/go/

Code documentation requirements based on Google Style Guide:

    Packages: Each package should begin with a comment describing its purpose.
    Functions: Document exported functions. Comments should start with the function name and briefly describe its behavior.
    Types: Add comments to exported structures and interfaces, describing their purpose.
    Variables and Constants: All exported variables and constants must be documented.
    Methods: Document exported methods, indicating the type they belong to.
    Complex Code: Add explanatory comments for complex parts of the code.
    Style: Use single-line comments // and write in the third person.
    Clarity: Comments should explain why something is done, not how.
    Language: Use English for comments.

Expected Results:

    A server application with support for:
        Configuring IKEv2 VPN.
        Generating configuration files for iPhone (.mobileconfig).
        Generating configuration files for OpenVPN (.ovpn).
        Deleting all client data after completion and providing a new password to the client.
    Documentation:
        Server deployment instructions.
        Examples of API requests.
    A test server to verify functionality.


Frontend: 

Flutter. Must have a GUI that allows the setup of the VPNs on any remote server.

Docker: Must use docker-compose. Do logging of every service and function. Ensure that the WHOLE application is containerized.