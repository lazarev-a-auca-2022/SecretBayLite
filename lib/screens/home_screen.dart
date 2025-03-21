import 'package:flutter/material.dart';
import 'package:secretbay_lite/widgets/server_config_form.dart';
import 'package:secretbay_lite/widgets/progress_indicator.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  bool isConfiguring = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('SecretBay Lite'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            ServerConfigForm(
              onSubmit: (config) async {
                setState(() => isConfiguring = true);
                // TODO: Implement VPN configuration
                setState(() => isConfiguring = false);
              },
            ),
            if (isConfiguring) 
              const ConfigurationProgress(),
          ],
        ),
      ),
    );
  }
}
