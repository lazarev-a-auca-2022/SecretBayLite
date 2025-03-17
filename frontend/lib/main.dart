import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:secretbay/screens/home_screen.dart';
import 'package:secretbay/services/auth_service.dart';
import 'package:secretbay/services/vpn_service.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthService()),
        ChangeNotifierProvider(create: (_) => VPNService()),
      ],
      child: MaterialApp(
        title: 'SecretBay',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          primarySwatch: Colors.blue,
          visualDensity: VisualDensity.adaptivePlatformDensity,
        ),
        home: const HomeScreen(),
      ),
    );
  }
}