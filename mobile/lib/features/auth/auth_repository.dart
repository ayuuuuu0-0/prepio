import '../../core/api/api_client.dart';

/// AuthRepository handles login and registration.
class AuthRepository {
  AuthRepository(this._api);
  final ApiClient _api;

  Future<AuthResult> login(String email, String password) => _api.login(email, password);

  Future<AuthResult> register(String email, String username, String password) =>
      _api.register(email, username, password);
}
