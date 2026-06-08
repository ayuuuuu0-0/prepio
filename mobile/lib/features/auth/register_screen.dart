import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/game_button.dart';
import '../../widgets/game/speech_bubble.dart';
import 'auth_provider.dart';
import 'profile_provider.dart';

/// RegisterScreen creates a new Prepio account with game styling.
class RegisterScreen extends ConsumerStatefulWidget {
  const RegisterScreen({super.key});

  @override
  ConsumerState<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends ConsumerState<RegisterScreen> {
  final _email = TextEditingController();
  final _username = TextEditingController();
  final _password = TextEditingController();
  var _loading = false;
  String? _error;

  @override
  void dispose() {
    _email.dispose();
    _username.dispose();
    _password.dispose();
    super.dispose();
  }

  Future<void> _register() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final repo = ref.read(authRepositoryProvider);
      final result = await repo.register(_email.text.trim(), _username.text.trim(), _password.text);
      await ref.read(sessionStoreProvider).saveToken(result.accessToken);
      ref.read(authTokenProvider.notifier).state = result.accessToken;
      ref.read(apiClientProvider).token = result.accessToken;
      ref.invalidate(profileProvider);
      if (mounted) Navigator.of(context).pop();
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
              const CompanionHero(name: 'Pip', species: 'red_panda', size: 110),
              const SizedBox(height: 20),
              const SpeechBubble(text: 'Join the adventure! Pick your companion and become interview-ready.'),
              if (_error != null) ...[
                const SizedBox(height: 16),
                Text(_error!, style: const TextStyle(color: Colors.orange, fontWeight: FontWeight.w600)),
              ],
              const SizedBox(height: 24),
              TextField(controller: _email, decoration: const InputDecoration(labelText: 'Email')),
              const SizedBox(height: 12),
              TextField(controller: _username, decoration: const InputDecoration(labelText: 'Username')),
              const SizedBox(height: 12),
              TextField(controller: _password, obscureText: true, decoration: const InputDecoration(labelText: 'Password')),
              const SizedBox(height: 24),
              GameButton(
                label: 'Begin Adventure!',
                color: PrepioColors.gold,
                shadowColor: const Color(0xFFE5B000),
                onPressed: _register,
                loading: _loading,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
