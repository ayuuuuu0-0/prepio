import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';

/// SpeechBubble displays companion dialogue in a dark HUD-style bubble.
class SpeechBubble extends StatelessWidget {
  const SpeechBubble({super.key, required this.text, this.speakerName});

  final String text;
  final String? speakerName;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: PrepioColors.surface,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: Colors.white.withValues(alpha: 0.08)),
        boxShadow: const [BoxShadow(color: Colors.black45, blurRadius: 12, offset: Offset(0, 4))],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (speakerName != null && speakerName!.isNotEmpty)
            Text(
              speakerName!.toUpperCase(),
              style: GoogleFonts.jetBrainsMono(
                fontSize: 10,
                fontWeight: FontWeight.w600,
                letterSpacing: 1.5,
                color: PrepioColors.accent,
              ),
            ),
          Text(
            text,
            style: GoogleFonts.nunito(fontWeight: FontWeight.w500, fontSize: 14, color: PrepioColors.textBody),
          ),
        ],
      ),
    );
  }
}
