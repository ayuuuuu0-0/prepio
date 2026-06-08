/// AppConfig holds runtime configuration for the mobile client.
class AppConfig {
  /// API gateway base URL. Override for physical devices:
  /// `flutter run --dart-define=API_BASE_URL=http://192.168.x.x:8080`
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );
}
