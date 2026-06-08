import 'package:flutter/material.dart';
import '../../core/theme/design_tokens.dart';

/// GameBackground wraps screens with gradient worlds — never plain white.
class GameBackground extends StatelessWidget {
  const GameBackground({super.key, required this.child, this.variant = GameBgVariant.defaultBg});

  final Widget child;
  final GameBgVariant variant;

  @override
  Widget build(BuildContext context) {
    final gradient = switch (variant) {
      GameBgVariant.defaultBg => const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [PrepioColors.bgTop, Color(0xFFF0FFE8), PrepioColors.bgBottom],
        ),
      GameBgVariant.forest => const LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [Color(0xFFB8E986), Color(0xFF7EC850), Color(0xFF5A9E3A)],
        ),
      GameBgVariant.challenge => const LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [Color(0xFFE8F5FF), Color(0xFFC8E6FF)],
        ),
    };

    return Container(
      decoration: BoxDecoration(gradient: gradient),
      child: SafeArea(child: child),
    );
  }
}

enum GameBgVariant { defaultBg, forest, challenge }
