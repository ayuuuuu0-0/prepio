import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/game_button.dart';
import '../challenge/challenge_screen.dart';

/// JourneyScreen shows the Foundation Forest path — hero feature placeholder.
class JourneyScreen extends StatelessWidget {
  const JourneyScreen({super.key});

  @override
  Widget build(BuildContext context) {
    const nodes = [
      ('done', 'Arrays', '✓'),
      ('done', 'Strings', '✓'),
      ('current', 'Hashing', '⚡'),
      ('locked', 'Sorting', '🔒'),
      ('boss', 'Boss', '👑'),
    ];

    return Scaffold(
      body: GameBackground(
        variant: GameBgVariant.forest,
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            Text('🌲 Foundation Forest', textAlign: TextAlign.center, style: GoogleFonts.fredoka(fontSize: 24, fontWeight: FontWeight.w800, color: Colors.white)),
            Text('World 1 · Beginner Journey', textAlign: TextAlign.center, style: GoogleFonts.nunito(color: Colors.white70, fontWeight: FontWeight.w600)),
            const SizedBox(height: 32),
            ...nodes.asMap().entries.map((entry) {
              final i = entry.key;
              final (type, label, icon) = entry.value;
              final isCurrent = type == 'current';
              final isBoss = type == 'boss';
              final align = i.isEven ? Alignment.centerLeft : Alignment.centerRight;

              return Align(
                alignment: align,
                child: Padding(
                  padding: const EdgeInsets.symmetric(vertical: 12),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (isCurrent) ...[
                        const CompanionHero(size: 48),
                        const SizedBox(width: 8),
                      ],
                      Column(
                        children: [
                          Container(
                            width: isBoss ? 72 : 60,
                            height: isBoss ? 72 : 60,
                            decoration: BoxDecoration(
                              shape: BoxShape.circle,
                              color: _nodeColor(type),
                              boxShadow: isCurrent ? [BoxShadow(color: Colors.cyan.withValues(alpha: 0.5), blurRadius: 12)] : null,
                            ),
                            alignment: Alignment.center,
                            child: Text(icon, style: TextStyle(fontSize: isBoss ? 28 : 22, color: Colors.white)),
                          ),
                          const SizedBox(height: 6),
                          Text(label, style: GoogleFonts.fredoka(fontSize: 12, fontWeight: FontWeight.w700, color: Colors.white)),
                        ],
                      ),
                    ],
                  ),
                ),
              );
            }),
            const SizedBox(height: 24),
            GameButton(
              label: 'Play Current Node →',
              onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const ChallengeScreen())),
            ),
          ],
        ),
      ),
    );
  }

  Color _nodeColor(String type) => switch (type) {
        'done' => const Color(0xFF58CC02),
        'current' => const Color(0xFF1CB0F6),
        'boss' => const Color(0xFFFFC800),
        _ => Colors.grey,
      };
}
