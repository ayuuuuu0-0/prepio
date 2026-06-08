import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/game_button.dart';
import '../../widgets/game/speech_bubble.dart';
import 'auth_provider.dart';
import 'profile_provider.dart';
import 'register_screen.dart';

/// LoginScreen handles email/password authentication with career RPG styling.
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
      await ref.read(sessionStoreProvider).saveTokens(
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      );
      ref.read(authTokenProvider.notifier).state = result.accessToken;
      final api = ref.read(apiClientProvider);
      api.token = result.accessToken;
      api.refreshToken = result.refreshToken;
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
      backgroundColor: PrepioColors.bg,
      body: GameBackground(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('PREPIO', style: GoogleFonts.plusJakartaSans(fontSize: 24, fontWeight: FontWeight.w800, color: PrepioColors.accent)),
              const SizedBox(height: 16),
              const CompanionHero(name: 'Byte', species: 'capybara', size: 100),
              const SizedBox(height: 20),
              const SpeechBubble(
                speakerName: 'Byte',
                text: "Byte is waiting. Your streak won't hold itself.",
              ),
              if (_error != null) ...[
                const SizedBox(height: 16),
                Text(_error!, style: const TextStyle(color: PrepioColors.danger, fontWeight: FontWeight.w600)),
              ],
              const SizedBox(height: 24),
              Text('Continue prep', style: GoogleFonts.plusJakartaSans(fontSize: 22, fontWeight: FontWeight.w700, color: PrepioColors.textPrimary)),
              const SizedBox(height: 8),
              Text('Sign in to your account', style: GoogleFonts.nunito(color: PrepioColors.textDim)),
              const SizedBox(height: 20),
              TextField(controller: _email, decoration: const InputDecoration(labelText: 'Email')),
              const SizedBox(height: 12),
              TextField(controller: _password, obscureText: true, decoration: const InputDecoration(labelText: 'Password')),
              const SizedBox(height: 24),
              GameButton(label: 'Continue Prep', onPressed: _login, loading: _loading),
              const SizedBox(height: 16),
              TextButton(
                onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const RegisterScreen())),
                child: Text(
                  'No account? Join 12k engineers in prep',
                  style: GoogleFonts.nunito(color: PrepioColors.accent, fontWeight: FontWeight.w600),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
