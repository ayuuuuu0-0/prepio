import 'package:hive_flutter/hive_flutter.dart';

/// SessionStore persists auth tokens across app restarts.
class SessionStore {
  static const _boxName = 'session';
  static const _accessKey = 'access_token';
  static const _refreshKey = 'refresh_token';

  Future<Box<String>> _box() => Hive.openBox<String>(_boxName);

  /// saveTokens stores access and refresh tokens locally.
  Future<void> saveTokens({required String accessToken, required String refreshToken}) async {
    final box = await _box();
    await box.put(_accessKey, accessToken);
    await box.put(_refreshKey, refreshToken);
  }

  /// saveToken stores only the access token (legacy).
  Future<void> saveToken(String token) async {
    final box = await _box();
    await box.put(_accessKey, token);
  }

  /// loadToken returns the stored access token, if any.
  Future<String?> loadToken() async {
    final box = await _box();
    return box.get(_accessKey);
  }

  /// loadRefreshToken returns the stored refresh token, if any.
  Future<String?> loadRefreshToken() async {
    final box = await _box();
    return box.get(_refreshKey);
  }

  /// clearToken removes stored tokens.
  Future<void> clearToken() async {
    final box = await _box();
    await box.delete(_accessKey);
    await box.delete(_refreshKey);
  }
}
