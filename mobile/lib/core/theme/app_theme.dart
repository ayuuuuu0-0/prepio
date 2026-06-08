import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'design_tokens.dart';

/// AppTheme builds the dark career RPG Material theme.
class AppTheme {
  static ThemeData get dark {
    final display = GoogleFonts.plusJakartaSansTextTheme();
    final body = GoogleFonts.nunitoTextTheme();
    final mono = GoogleFonts.jetBrainsMonoTextTheme();

    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      scaffoldBackgroundColor: PrepioColors.bg,
      colorScheme: const ColorScheme.dark(
        primary: PrepioColors.accent,
        secondary: PrepioColors.xp,
        surface: PrepioColors.surface,
        onSurface: PrepioColors.textPrimary,
        error: PrepioColors.danger,
      ),
      textTheme: body.copyWith(
        headlineLarge: display.headlineLarge?.copyWith(
          fontWeight: FontWeight.w700,
          color: PrepioColors.textPrimary,
        ),
        headlineMedium: display.headlineMedium?.copyWith(
          fontWeight: FontWeight.w700,
          color: PrepioColors.textPrimary,
        ),
        titleLarge: display.titleLarge?.copyWith(
          fontWeight: FontWeight.w700,
          color: PrepioColors.textPrimary,
        ),
        bodyMedium: body.bodyMedium?.copyWith(color: PrepioColors.textBody),
        labelSmall: mono.labelSmall?.copyWith(color: PrepioColors.textDim),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: PrepioColors.accent,
          foregroundColor: Colors.white,
          elevation: 0,
          padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 32),
          shape: const StadiumBorder(),
          textStyle: GoogleFonts.plusJakartaSans(fontSize: 16, fontWeight: FontWeight.w700),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: PrepioColors.surface,
        labelStyle: GoogleFonts.nunito(color: PrepioColors.textMuted),
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: PrepioColors.border),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: PrepioColors.accent, width: 1.5),
        ),
      ),
    );
  }
}
