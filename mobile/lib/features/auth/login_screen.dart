import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/game_button.dart';
import '../../widgets/game/speech_bubble.dart';
import 'auth_provider.dart';
import 'profile_provider.dart';
import 'register_screen.dart';

/// LoginScreen handles email/password authentication with game styling.
class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _email = TextEditingController();
  final _password = TextEditingController();
  var _loading = false;
  String? _error;

  @override
  void dispose() {
    _email.dispose();
    _password.dispose();
    super.dispose();
  }

  Future<void> _login() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final repo = ref.read(authRepositoryProvider);
      final result = await repo.login(_email.text.trim(), _password.text);
      await ref.read(sessionStoreProvider).saveToken(result.accessToken);
      ref.read(authTokenProvider.notifier).state = result.accessToken;
      ref.read(apiClientProvider).token = result.accessToken;
      ref.invalidate(profileProvider);
    } catch (e) {
      setState(() => _error = e.toString());
    } finally {
      setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: GameBackground(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const CompanionHero(name: 'Byte', species: 'capybara', size: 110),
              const SizedBox(height: 20),
              const SpeechBubble(text: 'Welcome back! Ready to level up your career today?'),
              if (_error != null) ...[
                const SizedBox(height: 16),
                Text(_error!, style: const TextStyle(color: Colors.orange, fontWeight: FontWeight.w600)),
              ],
              const SizedBox(height: 24),
              TextField(controller: _email, decoration: const InputDecoration(labelText: 'Email')),
              const SizedBox(height: 12),
              TextField(controller: _password, obscureText: true, decoration: const InputDecoration(labelText: 'Password')),
              const SizedBox(height: 24),
              GameButton(label: "Let's Go!", onPressed: _login, loading: _loading),
              const SizedBox(height: 16),
              TextButton(
                onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const RegisterScreen())),
                child: Text('New adventurer? Start your journey', style: GoogleFonts.nunito(fontWeight: FontWeight.w700, color: const Color(0xFF1CB0F6))),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
