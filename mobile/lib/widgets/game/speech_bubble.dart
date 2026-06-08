import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// SpeechBubble displays companion dialogue.
class SpeechBubble extends StatelessWidget {
  const SpeechBubble({super.key, required this.text});

  final String text;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24),
        boxShadow: const [BoxShadow(color: Colors.black12, blurRadius: 8, offset: Offset(0, 4))],
      ),
      child: Text(
        text,
        style: GoogleFonts.nunito(fontWeight: FontWeight.w700, fontSize: 14, color: const Color(0xFF3C3C3C)),
      ),
    );
  }
}
