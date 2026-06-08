import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/api/api_client.dart';
import '../auth/auth_provider.dart';

/// DailyPaperProvider loads today's question paper.
final dailyPaperProvider = FutureProvider<DailyPaper>((ref) async {
  final api = ref.watch(apiClientProvider);
  final token = ref.watch(authTokenProvider);
  if (token == null || token.isEmpty) {
    throw StateError('not authenticated');
  }
  api.token = token;
  return api.getDailyPaper();
});
