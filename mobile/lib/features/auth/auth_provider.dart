import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/api/api_client.dart';
import '../../core/config/app_config.dart';
import '../../core/storage/session_store.dart';
import 'auth_repository.dart';

final sessionStoreProvider = Provider<SessionStore>((ref) => SessionStore());

final apiClientProvider = Provider<ApiClient>(
  (ref) => ApiClient(baseUrl: AppConfig.apiBaseUrl),
);

final authRepositoryProvider = Provider<AuthRepository>(
  (ref) => AuthRepository(ref.watch(apiClientProvider)),
);

final authTokenProvider = StateProvider<String?>((ref) => null);

final authBootstrapProvider = FutureProvider<bool>((ref) async {
  final store = ref.read(sessionStoreProvider);
  final api = ref.read(apiClientProvider);
  final token = await store.loadToken();
  final refresh = await store.loadRefreshToken();
  if (token != null && token.isNotEmpty) {
    api.token = token;
    api.refreshToken = refresh;
    ref.read(authTokenProvider.notifier).state = token;
    return true;
  }
  if (refresh != null && refresh.isNotEmpty) {
    api.refreshToken = refresh;
    if (await api.refreshAccessToken()) {
      await store.saveTokens(accessToken: api.token!, refreshToken: api.refreshToken!);
      ref.read(authTokenProvider.notifier).state = api.token;
      return true;
    }
  }
  return false;
});
