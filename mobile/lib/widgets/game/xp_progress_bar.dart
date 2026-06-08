import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';

/// XpProgressBar shows progress toward the next level.
class XpProgressBar extends StatelessWidget {
  const XpProgressBar({super.key, required this.level, required this.totalXp, required this.xpToNextLevel});

  final int level;
  final int totalXp;
  final int xpToNextLevel;

  @override
  Widget build(BuildContext context) {
    final levelStart = levelThresholds[level - 1];
    final levelEnd = levelThresholds.length > level ? levelThresholds[level] : levelStart + 1000;
    final xpInLevel = totalXp - levelStart;
    final xpForLevel = levelEnd - levelStart;
    final pct = xpForLevel > 0 ? (xpInLevel / xpForLevel).clamp(0.0, 1.0) : 1.0;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: PrepioColors.surface,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: PrepioColors.border),
      ),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Level $level', style: GoogleFonts.jetBrainsMono(fontWeight: FontWeight.w700, color: PrepioColors.xp)),
              Text('Level ${level + 1}', style: GoogleFonts.jetBrainsMono(fontWeight: FontWeight.w700, color: PrepioColors.accent)),
            ],
          ),
          const SizedBox(height: 8),
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: LinearProgressIndicator(
              value: pct,
              minHeight: 8,
              backgroundColor: PrepioColors.border,
              color: PrepioColors.xp,
            ),
          ),
          const SizedBox(height: 6),
          Text(
            xpToNextLevel > 0 ? '$xpToNextLevel XP to next level' : 'Max level reached',
            style: GoogleFonts.jetBrainsMono(fontSize: 11, color: PrepioColors.textDim),
          ),
        ],
      ),
    );
  }
}
