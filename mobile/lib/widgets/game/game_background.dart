import 'package:flutter/material.dart';
import '../../core/theme/design_tokens.dart';

/// GameBackground wraps screens with dark ambient game backgrounds.
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
          colors: [Color(0xFF141622), PrepioColors.bg, Color(0xFF0F1419)],
        ),
      GameBgVariant.forest => const LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [Color(0xFF0D1F14), Color(0xFF0F1A10), PrepioColors.bg],
        ),
      GameBgVariant.challenge => const LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [Color(0xFF0F1020), PrepioColors.bg],
        ),
    };

    return Container(
      decoration: BoxDecoration(color: PrepioColors.bg, gradient: gradient),
      child: SafeArea(child: child),
    );
  }
}

enum GameBgVariant { defaultBg, forest, challenge }
