import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;

class VPNService extends ChangeNotifier {
  // Updated _baseUrl to use current origin and avoid CORS issues.
  static final String _baseUrl = '${Uri.base.origin}/api';
  bool _isConfiguring = false;
  String? _error;

  bool get isConfiguring => _isConfiguring;
  String? get error => _error;

  Future<Map<String, dynamic>> configureVPN({
    required String serverIp,
    required String username,
    required String authMethod,
    required String authCredential,
    required String vpnType,
  }) async {
    _isConfiguring = true;
    _error = null;
    notifyListeners();

    try {
      final response = await http.post(
        Uri.parse('$_baseUrl/vpn/configure'),
        headers: {
          'Content-Type': 'application/json',
        },
        body: jsonEncode({
          'server_ip': serverIp,
          'username': username,
          'auth_method': authMethod,
          'auth_credential': authCredential,
          'vpn_type': vpnType,
        }),
      );

      if (response.statusCode == 200) {
        _isConfiguring = false;
        notifyListeners();
        return jsonDecode(response.body);
      } else {
        throw Exception('Failed to configure VPN: ${response.body}');
      }
    } catch (e) {
      _error = e.toString();
      _isConfiguring = false;
      notifyListeners();
      rethrow;
    }
  }

  Future<bool> checkHealth() async {
    try {
      final response = await http.get(Uri.parse('$_baseUrl/health'));
      return response.statusCode == 200;
    } catch (e) {
      return false;
    }
  }
}