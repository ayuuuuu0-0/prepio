import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/api/api_client.dart';
import 'auth_provider.dart';

/// profileProvider loads the authenticated user's full profile.
final profileProvider = FutureProvider<ProfileInfo>((ref) async {
  final api = ref.watch(apiClientProvider);
  final token = ref.watch(authTokenProvider);
  if (token == null || token.isEmpty) {
    throw StateError('not authenticated');
  }
  api.token = token;
  return api.getProfile();
});
