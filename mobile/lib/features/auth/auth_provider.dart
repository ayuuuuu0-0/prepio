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
  final token = await store.loadToken();
  if (token != null && token.isNotEmpty) {
    ref.read(authTokenProvider.notifier).state = token;
    ref.read(apiClientProvider).token = token;
    return true;
  }
  return false;
});
