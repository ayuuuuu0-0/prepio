import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/api/api_client.dart';
import '../auth/auth_provider.dart';

/// ProgressProvider loads XP and gem balance.
final progressProvider = FutureProvider<ProgressInfo>((ref) async {
  final api = ref.watch(apiClientProvider);
  final token = ref.watch(authTokenProvider);
  if (token == null || token.isEmpty) {
    throw StateError('not authenticated');
  }
  api.token = token;
  return api.getProgress();
});
