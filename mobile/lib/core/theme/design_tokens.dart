import 'package:flutter/material.dart';

/// PrepioColors holds the career RPG dark palette.
class PrepioColors {
  static const bg = Color(0xFF0F1117);
  static const surface = Color(0xFF1A1D27);
  static const raised = Color(0xFF242836);
  static const border = Color(0xFF2E3347);
  static const accent = Color(0xFF7C6EF5);
  static const streak = Color(0xFFFF6B35);
  static const gems = Color(0xFF34D399);
  static const xp = Color(0xFF60A5FA);
  static const gold = Color(0xFFF5B942);
  static const success = Color(0xFF34D399);
  static const warning = Color(0xFFF5B942);
  static const danger = Color(0xFFF87171);
  static const textPrimary = Color(0xFFE8EAED);
  static const textMuted = Color(0xFF8B92A8);
  static const textDim = Color(0xFF4A5068);
  static const textBody = Color(0xFFC8CCDA);
}

/// LevelThresholds mirrors config/levels.go cumulative XP per level.
const levelThresholds = [0, 100, 250, 500, 800, 1200, 1700, 2300, 3000, 3800];

/// CompanyColors maps target companies to ring colors.
const companyRingColors = <String, Color>{
  'google': Color(0xFF4285F4),
  'amazon': Color(0xFFFF9900),
  'meta': Color(0xFF7C6EF5),
  'uber': Color(0xFF34D399),
  'atlassian': Color(0xFF0052CC),
  'netflix': Color(0xFFE50914),
};

/// RoundTypeColors maps question round types to accent colors.
const roundTypeColors = <String, Color>{
  'dsa': Color(0xFF7C6EF5),
  'system_design': Color(0xFF60A5FA),
  'lld': Color(0xFF34D399),
  'behavioral': Color(0xFFFF6B35),
  'aptitude': Color(0xFFF5B942),
  'fundamentals': Color(0xFFA99EFA),
};

/// CompanionVisual holds emoji and glow for each companion species.
class CompanionVisual {
  const CompanionVisual({required this.emoji, required this.glow});
  final String emoji;
  final Color glow;
}

/// companionFor returns visual identity by name or species.
CompanionVisual companionFor({String? name, String? species}) {
  final key = (name ?? species ?? 'byte').toLowerCase();
  const map = <String, CompanionVisual>{
    'byte': CompanionVisual(emoji: '🦫', glow: Color(0xFF7CB342)),
    'capybara': CompanionVisual(emoji: '🦫', glow: Color(0xFF7CB342)),
    'pip': CompanionVisual(emoji: '🐼', glow: Color(0xFFFF6B35)),
    'red_panda': CompanionVisual(emoji: '🐼', glow: Color(0xFFFF6B35)),
    'nova': CompanionVisual(emoji: '🦔', glow: Color(0xFF7C6EF5)),
    'pangolin': CompanionVisual(emoji: '🦔', glow: Color(0xFF7C6EF5)),
    'kodo': CompanionVisual(emoji: '🦎', glow: Color(0xFFFF6B9D)),
    'axolotl': CompanionVisual(emoji: '🦎', glow: Color(0xFFFF6B9D)),
    'zara': CompanionVisual(emoji: '🐆', glow: Color(0xFF60A5FA)),
    'snow_leopard': CompanionVisual(emoji: '🐆', glow: Color(0xFF60A5FA)),
  };
  return map[key] ?? map['byte']!;
}
