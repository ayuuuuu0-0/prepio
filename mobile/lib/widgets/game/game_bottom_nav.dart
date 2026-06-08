import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';

/// GameBottomNav provides 6-tab navigation matching the web app.
class GameBottomNav extends StatelessWidget {
  const GameBottomNav({super.key, required this.currentIndex, required this.onTap});

  final int currentIndex;
  final ValueChanged<int> onTap;

  static const tabs = [
    ('⊞', 'Home'),
    ('◈', 'Journey'),
    ('▶', 'Play'),
    ('⬡', 'League'),
    ('◇', 'Quests'),
    ('◉', 'Profile'),
  ];

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: PrepioColors.bg.withValues(alpha: 0.95),
        border: const Border(top: BorderSide(color: PrepioColors.border)),
      ),
      child: SafeArea(
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: List.generate(tabs.length, (i) {
            final (icon, label) = tabs[i];
            final active = i == currentIndex;
            final dimmed = i == 3 || i == 4;
            return GestureDetector(
              onTap: () => onTap(i),
              child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 6, horizontal: 4),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      icon,
                      style: TextStyle(
                        fontSize: 16,
                        color: active ? PrepioColors.accent : PrepioColors.textDim,
                      ),
                    ),
                    Text(
                      label,
                      style: GoogleFonts.plusJakartaSans(
                        fontSize: 9,
                        fontWeight: FontWeight.w600,
                        color: active ? PrepioColors.accent : PrepioColors.textDim,
                      ),
                    ),
                    if (dimmed) const SizedBox(height: 2),
                  ],
                ),
              ),
            );
          }),
        ),
      ),
    );
  }
}
