import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';

/// QuestCard displays a daily quest with progress bar and rewards.
class QuestCard extends StatelessWidget {
  const QuestCard({
    super.key,
    required this.title,
    required this.icon,
    required this.progress,
    required this.target,
    required this.completed,
    required this.rewardXp,
    required this.rewardGems,
    this.comingSoon = false,
  });

  final String title;
  final String icon;
  final int progress;
  final int target;
  final bool completed;
  final int rewardXp;
  final int rewardGems;
  final bool comingSoon;

  @override
  Widget build(BuildContext context) {
    if (comingSoon) {
      return Container(
        margin: const EdgeInsets.only(bottom: 10),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: PrepioColors.surface,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: PrepioColors.border, style: BorderStyle.solid),
        ),
        child: Row(
          children: [
            Text(icon, style: const TextStyle(fontSize: 22, color: PrepioColors.textDim)),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: GoogleFonts.plusJakartaSans(fontWeight: FontWeight.w700, color: PrepioColors.textDim)),
                  Text('Coming soon', style: GoogleFonts.jetBrainsMono(fontSize: 11, color: PrepioColors.textDim)),
                ],
              ),
            ),
          ],
        ),
      );
    }

    final pct = (progress / target).clamp(0.0, 1.0);

    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: completed ? PrepioColors.success.withValues(alpha: 0.08) : PrepioColors.surface,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: completed ? PrepioColors.success.withValues(alpha: 0.3) : PrepioColors.border,
        ),
      ),
      child: Row(
        children: [
          Text(completed ? '✅' : icon, style: const TextStyle(fontSize: 22)),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: GoogleFonts.plusJakartaSans(
                    fontWeight: FontWeight.w700,
                    decoration: completed ? TextDecoration.lineThrough : null,
                    color: completed ? PrepioColors.success : PrepioColors.textPrimary,
                  ),
                ),
                const SizedBox(height: 8),
                ClipRRect(
                  borderRadius: BorderRadius.circular(8),
                  child: LinearProgressIndicator(
                    value: pct,
                    minHeight: 6,
                    backgroundColor: PrepioColors.border,
                    color: completed ? PrepioColors.success : PrepioColors.accent,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '$progress/$target · ⚡ $rewardXp XP · 💎 $rewardGems',
                  style: GoogleFonts.jetBrainsMono(fontSize: 11, color: PrepioColors.textDim),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
