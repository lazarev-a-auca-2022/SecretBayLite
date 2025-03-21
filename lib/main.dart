import 'package:flutter/material.dart';
import 'package:secretbay_lite/screens/home_screen.dart';
import 'package:secretbay_lite/theme/app_theme.dart';

void main() {
  runApp(const SecretBayApp());
}

class SecretBayApp extends StatelessWidget {
  const SecretBayApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'SecretBay Lite',
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      home: const HomeScreen(),
    );
  }
}
