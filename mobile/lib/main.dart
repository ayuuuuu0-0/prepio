import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'core/theme/app_theme.dart';
import 'features/auth/auth_provider.dart';
import 'features/auth/login_screen.dart';
import 'features/auth/profile_provider.dart';
import 'features/dashboard/dashboard_screen.dart';
import 'features/onboarding/onboarding_screen.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Hive.initFlutter();
  await Hive.openBox('pending_answers');
  runApp(const ProviderScope(child: PrepioApp()));
}

class PrepioApp extends ConsumerWidget {
  const PrepioApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final bootstrap = ref.watch(authBootstrapProvider);
    final token = ref.watch(authTokenProvider);

    return MaterialApp(
      title: 'Prepio',
      theme: AppTheme.light,
      debugShowCheckedModeBanner: false,
      home: bootstrap.when(
        loading: () => const Scaffold(body: Center(child: CircularProgressIndicator())),
        error: (_, __) => const LoginScreen(),
        data: (_) {
          if (token == null || token.isEmpty) return const LoginScreen();
          return const _AuthenticatedRoot();
        },
      ),
    );
  }
}

class _AuthenticatedRoot extends ConsumerWidget {
  const _AuthenticatedRoot();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final profile = ref.watch(profileProvider);
    return profile.when(
      loading: () => const Scaffold(body: Center(child: CircularProgressIndicator())),
      error: (_, __) => const LoginScreen(),
      data: (p) => p.onboardingCompleted ? const DashboardScreen() : const OnboardingScreen(),
    );
  }
}
