import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'features/auth/login_screen.dart';
import 'features/home/home_screen.dart';
import 'features/auth/auth_provider.dart';

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
    final token = ref.watch(authTokenProvider);

    return MaterialApp(
      title: 'Prepio',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF059669)),
        useMaterial3: true,
      ),
      home: token != null && token.isNotEmpty ? const HomeScreen() : const LoginScreen(),
    );
  }
}
