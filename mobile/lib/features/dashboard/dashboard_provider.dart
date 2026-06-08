import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/api/api_client.dart';
import '../auth/auth_provider.dart';

/// dashboardProvider loads the aggregated home dashboard.
final dashboardProvider = FutureProvider<DashboardHome>((ref) async {
  final api = ref.watch(apiClientProvider);
  final token = ref.watch(authTokenProvider);
  if (token == null || token.isEmpty) {
    throw StateError('not authenticated');
  }
  api.token = token;
  return api.getDashboardHome();
});
