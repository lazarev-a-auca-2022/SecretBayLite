import 'dart:convert';
import 'package:http/http.dart' as http;

class VPNService {
  static const String _baseUrl = 'http://localhost:8080/api';

  Future<Map<String, dynamic>> configureVPN({
    required String serverIP,
    required String username,
    required String authMethod,
    required String authCredential,
    required String vpnType,
  }) async {
    try {
      final response = await http.post(
        Uri.parse('$_baseUrl/vpn/configure'),
        headers: {
          'Content-Type': 'application/json',
        },
        body: jsonEncode({
          'server_ip': serverIP,
          'username': username,
          'auth_method': authMethod,
          'auth_credential': authCredential,
          'vpn_type': vpnType,
        }),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Failed to configure VPN: ${response.body}');
      }
    } catch (e) {
      throw Exception('Network error: $e');
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