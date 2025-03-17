import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/vpn_service.dart';
import '../widgets/config_form.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('SecretBay VPN Configuration'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () async {
              final vpnService = context.read<VPNService>();
              final isHealthy = await vpnService.checkHealth();
              if (context.mounted) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      isHealthy ? 'Server is healthy' : 'Server is not responding',
                    ),
                    backgroundColor: isHealthy ? Colors.green : Colors.red,
                  ),
                );
              }
            },
          ),
        ],
      ),
      body: const Padding(
        padding: EdgeInsets.all(16.0),
        child: ConfigForm(),
      ),
    );
  }
}