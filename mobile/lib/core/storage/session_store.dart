import 'package:hive_flutter/hive_flutter.dart';

/// SessionStore persists the access token across app restarts.
class SessionStore {
  static const _boxName = 'session';
  static const _tokenKey = 'access_token';

  Future<Box<String>> _box() => Hive.openBox<String>(_boxName);

  /// saveToken stores the access token locally.
  Future<void> saveToken(String token) async {
    final box = await _box();
    await box.put(_tokenKey, token);
  }

  /// loadToken returns the stored access token, if any.
  Future<String?> loadToken() async {
    final box = await _box();
    return box.get(_tokenKey);
  }

  /// clearToken removes the stored access token.
  Future<void> clearToken() async {
    final box = await _box();
    await box.delete(_tokenKey);
  }
}
