import 'package:flutter_test/flutter_test.dart';
import 'package:prepio/core/api/api_client.dart';

void main() {
  test('AuthResult parses gateway response', () {
    final result = AuthResult.fromJson({
      'access_token': 'tok',
      'refresh_token': 'ref',
      'user': {'id': '1', 'username': 'demo', 'email': 'a@test.com'},
    });
    expect(result.accessToken, 'tok');
    expect(result.username, 'demo');
  });
}
