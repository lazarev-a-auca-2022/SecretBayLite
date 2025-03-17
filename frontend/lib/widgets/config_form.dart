import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:file_picker/file_picker.dart';
import '../services/vpn_service.dart';

class ConfigForm extends StatefulWidget {
  const ConfigForm({super.key});

  @override
  State<ConfigForm> createState() => _ConfigFormState();
}

class _ConfigFormState extends State<ConfigForm> {
  final _formKey = GlobalKey<FormState>();
  final _serverIpController = TextEditingController();
  final _usernameController = TextEditingController();
  final _authCredentialController = TextEditingController();
  String _authMethod = 'password';
  String _vpnType = 'openvpn';
  bool _isLoading = false;

  @override
  void dispose() {
    _serverIpController.dispose();
    _usernameController.dispose();
    _authCredentialController.dispose();
    super.dispose();
  }

  Future<void> _uploadKeyFile() async {
    final result = await FilePicker.platform.pickFiles(
      type: FileType.custom,
      allowedExtensions: ['pem', 'key'],
    );

    if (result != null) {
      _authCredentialController.text = result.files.first.path ?? '';
    }
  }

  Future<void> _submitForm() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    try {
      final vpnService = context.read<VPNService>();
      final result = await vpnService.configureVPN(
        serverIP: _serverIpController.text,
        username: _usernameController.text,
        authMethod: _authMethod,
        authCredential: _authCredentialController.text,
        vpnType: _vpnType,
      );

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('VPN configured successfully!')),
        );
        // Show configuration details in a dialog
        _showConfigurationDialog(result);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    } finally {
      setState(() => _isLoading = false);
    }
  }

  void _showConfigurationDialog(Map<String, dynamic> config) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('VPN Configuration'),
        content: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text('Configuration:', style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Text(config['config'] ?? 'No configuration available'),
              const SizedBox(height: 16),
              const Text('New Password:', style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Text(config['new_password'] ?? 'No password available'),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: ListView(
        children: [
          TextFormField(
            controller: _serverIpController,
            decoration: const InputDecoration(
              labelText: 'Server IP',
              hintText: 'Enter server IP address',
            ),
            validator: (value) {
              if (value?.isEmpty ?? true) {
                return 'Please enter server IP';
              }
              return null;
            },
          ),
          const SizedBox(height: 16),
          TextFormField(
            controller: _usernameController,
            decoration: const InputDecoration(
              labelText: 'Username',
              hintText: 'Enter username',
            ),
            validator: (value) {
              if (value?.isEmpty ?? true) {
                return 'Please enter username';
              }
              return null;
            },
          ),
          const SizedBox(height: 16),
          DropdownButtonFormField<String>(
            value: _authMethod,
            decoration: const InputDecoration(
              labelText: 'Authentication Method',
            ),
            items: const [
              DropdownMenuItem(value: 'password', child: Text('Password')),
              DropdownMenuItem(value: 'key', child: Text('SSH Key')),
            ],
            onChanged: (value) {
              setState(() {
                _authMethod = value!;
                _authCredentialController.clear();
              });
            },
          ),
          const SizedBox(height: 16),
          if (_authMethod == 'password')
            TextFormField(
              controller: _authCredentialController,
              decoration: const InputDecoration(
                labelText: 'Password',
                hintText: 'Enter password',
              ),
              obscureText: true,
              validator: (value) {
                if (value?.isEmpty ?? true) {
                  return 'Please enter password';
                }
                return null;
              },
            )
          else
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: _authCredentialController,
                    decoration: const InputDecoration(
                      labelText: 'SSH Key File',
                      hintText: 'Select SSH key file',
                    ),
                    readOnly: true,
                    validator: (value) {
                      if (value?.isEmpty ?? true) {
                        return 'Please select SSH key file';
                      }
                      return null;
                    },
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.upload_file),
                  onPressed: _uploadKeyFile,
                ),
              ],
            ),
          const SizedBox(height: 16),
          DropdownButtonFormField<String>(
            value: _vpnType,
            decoration: const InputDecoration(
              labelText: 'VPN Type',
            ),
            items: const [
              DropdownMenuItem(value: 'openvpn', child: Text('OpenVPN')),
              DropdownMenuItem(value: 'ios', child: Text('IKEv2 (iOS)')),
            ],
            onChanged: (value) {
              setState(() => _vpnType = value!);
            },
          ),
          const SizedBox(height: 24),
          ElevatedButton(
            onPressed: _isLoading ? null : _submitForm,
            child: _isLoading
                ? const CircularProgressIndicator()
                : const Text('Configure VPN'),
          ),
        ],
      ),
    );
  }
}