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
  });

  final String title;
  final String icon;
  final int progress;
  final int target;
  final bool completed;
  final int rewardXp;
  final int rewardGems;

  @override
  Widget build(BuildContext context) {
    final pct = (progress / target).clamp(0.0, 1.0);

    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: completed ? const Color(0xFFE8F8D8) : Colors.white,
        borderRadius: BorderRadius.circular(24),
        border: completed ? Border.all(color: PrepioColors.green, width: 2) : null,
        boxShadow: const [BoxShadow(color: Colors.black12, blurRadius: 6, offset: Offset(0, 3))],
      ),
      child: Row(
        children: [
          Text(completed ? '✅' : icon, style: const TextStyle(fontSize: 28)),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: GoogleFonts.fredoka(
                    fontWeight: FontWeight.w700,
                    decoration: completed ? TextDecoration.lineThrough : null,
                    color: completed ? PrepioColors.greenDark : PrepioColors.text,
                  ),
                ),
                const SizedBox(height: 8),
                ClipRRect(
                  borderRadius: BorderRadius.circular(8),
                  child: LinearProgressIndicator(
                    value: pct,
                    minHeight: 10,
                    backgroundColor: const Color(0xFFE5E5E5),
                    color: PrepioColors.green,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '$progress/$target · ⚡$rewardXp · 💎$rewardGems',
                  style: GoogleFonts.nunito(fontSize: 11, color: PrepioColors.textMuted),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
