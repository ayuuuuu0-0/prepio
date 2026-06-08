import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';

/// GameButtonVariant selects the button color scheme.
enum GameButtonVariant { primary, secondary, gold }

/// GameButton is a large pill CTA with game-style press feedback.
class GameButton extends StatefulWidget {
  const GameButton({
    super.key,
    required this.label,
    required this.onPressed,
    this.variant = GameButtonVariant.primary,
    this.loading = false,
  });

  final String label;
  final VoidCallback? onPressed;
  final GameButtonVariant variant;
  final bool loading;

  @override
  State<GameButton> createState() => _GameButtonState();
}

class _GameButtonState extends State<GameButton> {
  var _pressed = false;

  (Color bg, Color shadow, Color text) get _colors => switch (widget.variant) {
        GameButtonVariant.primary => (PrepioColors.accent, const Color(0xFF5B50D4), Colors.white),
        GameButtonVariant.secondary => (PrepioColors.raised, PrepioColors.surface, PrepioColors.textBody),
        GameButtonVariant.gold => (PrepioColors.gold, const Color(0xFFC8941A), const Color(0xFF1A1500)),
      };

  @override
  Widget build(BuildContext context) {
    final (bg, shadow, text) = _colors;

    return GestureDetector(
      onTapDown: widget.onPressed == null ? null : (_) => setState(() => _pressed = true),
      onTapUp: widget.onPressed == null ? null : (_) => setState(() => _pressed = false),
      onTapCancel: () => setState(() => _pressed = false),
      onTap: widget.loading ? null : widget.onPressed,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 100),
        transform: Matrix4.translationValues(0, _pressed ? 3 : 0, 0),
        padding: const EdgeInsets.symmetric(vertical: 16),
        decoration: BoxDecoration(
          color: bg,
          borderRadius: BorderRadius.circular(999),
          border: widget.variant == GameButtonVariant.secondary
              ? Border.all(color: PrepioColors.border)
              : null,
          boxShadow: [
            BoxShadow(
              color: shadow,
              offset: Offset(0, _pressed ? 2 : 4),
              blurRadius: 0,
            ),
          ],
        ),
        alignment: Alignment.center,
        child: widget.loading
            ? SizedBox(
                height: 22,
                width: 22,
                child: CircularProgressIndicator(color: text, strokeWidth: 2),
              )
            : Text(
                widget.label,
                style: GoogleFonts.plusJakartaSans(fontSize: 16, fontWeight: FontWeight.w700, color: text),
              ),
      ),
    );
  }
}
