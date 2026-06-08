import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// ReadinessRing shows company readiness as a circular progress indicator.
class ReadinessRing extends StatelessWidget {
  const ReadinessRing({super.key, required this.company, required this.score, required this.color});

  final String company;
  final int score;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SizedBox(
          width: 80,
          height: 80,
          child: Stack(
            alignment: Alignment.center,
            children: [
              CircularProgressIndicator(
                value: score / 100,
                strokeWidth: 8,
                backgroundColor: const Color(0xFFE8E8E8),
                color: color,
              ),
              Text(
                '$score%',
                style: GoogleFonts.fredoka(fontWeight: FontWeight.w700, fontSize: 16, color: color),
              ),
            ],
          ),
        ),
        const SizedBox(height: 6),
        Text(
          company[0].toUpperCase() + company.substring(1),
          style: GoogleFonts.fredoka(fontSize: 11, fontWeight: FontWeight.w700),
        ),
      ],
    );
  }
}
