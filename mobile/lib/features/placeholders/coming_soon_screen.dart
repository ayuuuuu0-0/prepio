import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/speech_bubble.dart';

/// ComingSoonScreen is a placeholder for features launching in a future phase.
class ComingSoonScreen extends StatelessWidget {
  const ComingSoonScreen({super.key, required this.title, required this.icon, required this.message});
  final String title;
  final String icon;
  final String message;

  @override
  Widget build(BuildContext context) {
    return GameBackground(
      child: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(icon, style: const TextStyle(fontSize: 48)),
              Text(title, style: GoogleFonts.plusJakartaSans(fontSize: 24, fontWeight: FontWeight.w800, color: PrepioColors.textPrimary)),
              const SizedBox(height: 16),
              const CompanionHero(size: 80),
              const SizedBox(height: 16),
              SpeechBubble(text: message),
            ],
          ),
        ),
      ),
    );
  }
}
