import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../auth/auth_provider.dart';
import '../auth/profile_provider.dart';

/// ProfileScreen shows the authenticated user's profile summary.
class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final profile = ref.watch(profileProvider);

    return GameBackground(
      child: profile.when(
        loading: () => const Center(child: CircularProgressIndicator(color: PrepioColors.accent)),
        error: (e, _) => Center(child: Text('$e')),
        data: (p) => ListView(
          padding: const EdgeInsets.all(24),
          children: [
            const SizedBox(height: 24),
            Center(child: CompanionHero(name: p.companion?.name, species: p.companion?.species, size: 110)),
            const SizedBox(height: 16),
            Text(p.companion?.name ?? 'Adventurer', textAlign: TextAlign.center, style: GoogleFonts.plusJakartaSans(fontSize: 24, fontWeight: FontWeight.w800, color: PrepioColors.textPrimary)),
            Text('Target: ${p.targetCompanies.join(', ')}', textAlign: TextAlign.center, style: GoogleFonts.nunito(color: PrepioColors.textMuted)),
            const SizedBox(height: 32),
            TextButton(
              onPressed: () async {
                await ref.read(sessionStoreProvider).clearToken();
                ref.read(authTokenProvider.notifier).state = null;
              },
              child: const Text('Sign out'),
            ),
          ],
        ),
      ),
    );
  }
}
