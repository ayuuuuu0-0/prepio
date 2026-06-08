import 'package:flutter/material.dart';

/// PrepioColors holds the game palette — never corporate grey/white.
class PrepioColors {
  static const green = Color(0xFF58CC02);
  static const greenDark = Color(0xFF46A302);
  static const purple = Color(0xFF7B5CFF);
  static const orange = Color(0xFFFF9600);
  static const blue = Color(0xFF1CB0F6);
  static const gold = Color(0xFFFFC800);
  static const pink = Color(0xFFFF6B9D);
  static const text = Color(0xFF3C3C3C);
  static const textMuted = Color(0xFF777777);
  static const bgTop = Color(0xFFE8F5D8);
  static const bgBottom = Color(0xFFC8E6FF);
}

/// CompanyColors maps target companies to ring colors.
const companyRingColors = <String, Color>{
  'google': Color(0xFF1CB0F6),
  'amazon': Color(0xFFFF9600),
  'meta': Color(0xFF7B5CFF),
  'uber': Color(0xFF58CC02),
  'atlassian': Color(0xFF0052CC),
};

/// CompanionVisual holds emoji and gradient for each companion species.
class CompanionVisual {
  const CompanionVisual({required this.emoji, required this.colors});
  final String emoji;
  final List<Color> colors;
}

/// companionFor returns visual identity by name or species.
CompanionVisual companionFor({String? name, String? species}) {
  final key = (name ?? species ?? 'byte').toLowerCase();
  const map = <String, CompanionVisual>{
    'byte': CompanionVisual(emoji: '🦫', colors: [Color(0xFF8B5E3C), Color(0xFFD4A574), Color(0xFF7CB342)]),
    'capybara': CompanionVisual(emoji: '🦫', colors: [Color(0xFF8B5E3C), Color(0xFFD4A574), Color(0xFF7CB342)]),
    'pip': CompanionVisual(emoji: '🐼', colors: [Color(0xFFFF6B35), Color(0xFFFFB84D), Color(0xFFE63946)]),
    'red_panda': CompanionVisual(emoji: '🐼', colors: [Color(0xFFFF6B35), Color(0xFFFFB84D), Color(0xFFE63946)]),
    'nova': CompanionVisual(emoji: '🦔', colors: [Color(0xFF7B5CFF), Color(0xFF9B59B6), Color(0xFFFFC800)]),
    'pangolin': CompanionVisual(emoji: '🦔', colors: [Color(0xFF7B5CFF), Color(0xFF9B59B6), Color(0xFFFFC800)]),
    'kodo': CompanionVisual(emoji: '🦎', colors: [Color(0xFFFF85C0), Color(0xFFFFB4E6), Color(0xFF87CEEB)]),
    'axolotl': CompanionVisual(emoji: '🦎', colors: [Color(0xFFFF85C0), Color(0xFFFFB4E6), Color(0xFF87CEEB)]),
    'zara': CompanionVisual(emoji: '🐆', colors: [Color(0xFF94A3B8), Color(0xFFBAE6FD), Color(0xFFFFFFFF)]),
    'snow_leopard': CompanionVisual(emoji: '🐆', colors: [Color(0xFF94A3B8), Color(0xFFBAE6FD), Color(0xFFFFFFFF)]),
  };
  return map[key] ?? map['byte']!;
}
