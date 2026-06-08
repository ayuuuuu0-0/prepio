import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'design_tokens.dart';

/// AppTheme builds the playful game Material theme.
class AppTheme {
  static ThemeData get light {
    final display = GoogleFonts.fredokaTextTheme();
    final body = GoogleFonts.nunitoTextTheme();

    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: PrepioColors.green,
        primary: PrepioColors.green,
        secondary: PrepioColors.blue,
        surface: Colors.white,
      ),
      textTheme: body.copyWith(
        headlineLarge: display.headlineLarge?.copyWith(
          fontWeight: FontWeight.w700,
          color: PrepioColors.text,
        ),
        headlineMedium: display.headlineMedium?.copyWith(
          fontWeight: FontWeight.w700,
          color: PrepioColors.text,
        ),
        titleLarge: display.titleLarge?.copyWith(fontWeight: FontWeight.w700),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: PrepioColors.green,
          foregroundColor: Colors.white,
          elevation: 4,
          shadowColor: PrepioColors.greenDark,
          padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 32),
          shape: const StadiumBorder(),
          textStyle: GoogleFonts.fredoka(fontSize: 18, fontWeight: FontWeight.w700),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: Colors.white,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(20)),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(20),
          borderSide: const BorderSide(color: Color(0xFFE5E5E5), width: 2),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(20),
          borderSide: const BorderSide(color: PrepioColors.green, width: 2),
        ),
      ),
    );
  }
}
