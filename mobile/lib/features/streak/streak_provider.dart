import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/api/api_client.dart';
import '../auth/auth_provider.dart';

/// StreakProvider loads streak data from the API.
final streakProvider = FutureProvider<StreakInfo>((ref) async {
  final api = ref.watch(apiClientProvider);
  final token = ref.watch(authTokenProvider);
  if (token == null || token.isEmpty) {
    throw StateError('not authenticated');
  }
  api.token = token;
  return api.getStreak();
});
