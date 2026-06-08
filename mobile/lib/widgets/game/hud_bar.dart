import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/api/api_client.dart';
import '../../core/theme/design_tokens.dart';

/// HudBar combines companion badge, streak, gems, and XP progress in one row.
class HudBar extends StatelessWidget {
  const HudBar({super.key, required this.home});

  final DashboardHome home;

  @override
  Widget build(BuildContext context) {
    final progress = home.progress;
    final levelStart = levelThresholds[progress.currentLevel - 1];
    final levelEnd = levelThresholds.length > progress.currentLevel
        ? levelThresholds[progress.currentLevel]
        : levelStart + 1000;
    final xpInLevel = progress.totalXp - levelStart;
    final xpForLevel = levelEnd - levelStart;
    final pct = xpForLevel > 0 ? (xpInLevel / xpForLevel).clamp(0.0, 1.0) : 1.0;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: PrepioColors.surface,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: PrepioColors.border),
      ),
      child: Row(
        children: [
          Text(_companionEmoji(home.companion?.species), style: const TextStyle(fontSize: 20)),
          const SizedBox(width: 8),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                (home.companion?.name ?? 'Byte').toUpperCase(),
                style: GoogleFonts.jetBrainsMono(fontSize: 9, fontWeight: FontWeight.w700, color: PrepioColors.accent),
              ),
              Text(
                'Lv.${progress.currentLevel}',
                style: GoogleFonts.jetBrainsMono(fontSize: 12, fontWeight: FontWeight.w700, color: PrepioColors.xp),
              ),
            ],
          ),
          _divider(),
          Text('🔥', style: GoogleFonts.jetBrainsMono(fontSize: 14, fontWeight: FontWeight.w700, color: PrepioColors.streak)),
          Text('${home.streak.currentStreak}', style: GoogleFonts.jetBrainsMono(fontSize: 14, fontWeight: FontWeight.w700, color: PrepioColors.streak)),
          const SizedBox(width: 8),
          Text('💎', style: const TextStyle(fontSize: 14)),
          Text('${progress.gemBalance}', style: GoogleFonts.jetBrainsMono(fontSize: 14, fontWeight: FontWeight.w700, color: PrepioColors.gems)),
          _divider(),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text('XP', style: GoogleFonts.jetBrainsMono(fontSize: 9, color: PrepioColors.textDim)),
                    Text('${progress.xpToNextLevel} to next', style: GoogleFonts.jetBrainsMono(fontSize: 9, color: PrepioColors.xp)),
                  ],
                ),
                const SizedBox(height: 4),
                ClipRRect(
                  borderRadius: BorderRadius.circular(8),
                  child: LinearProgressIndicator(
                    value: pct,
                    minHeight: 6,
                    backgroundColor: PrepioColors.border,
                    color: PrepioColors.xp,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _divider() => Container(
        margin: const EdgeInsets.symmetric(horizontal: 8),
        width: 1,
        height: 28,
        color: PrepioColors.border,
      );

  String _companionEmoji(String? species) {
    const map = {
      'capybara': '🦫',
      'red_panda': '🐼',
      'pangolin': '🦔',
      'axolotl': '🦎',
      'snow_leopard': '🐆',
      'owl': '🦉',
    };
    return map[species ?? ''] ?? '🦫';
  }
}
