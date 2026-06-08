import 'package:flutter/material.dart';
import '../challenge/challenge_screen.dart';
import '../dashboard/dashboard_screen.dart';
import '../journey/journey_screen.dart';
import '../placeholders/coming_soon_screen.dart';
import '../profile/profile_screen.dart';
import '../../widgets/game/game_bottom_nav.dart';

/// AppShell hosts the main 6-tab navigation shell.
class AppShell extends StatefulWidget {
  const AppShell({super.key});

  @override
  State<AppShell> createState() => _AppShellState();
}

class _AppShellState extends State<AppShell> {
  var _index = 0;

  @override
  Widget build(BuildContext context) {
    final screens = [
      const DashboardScreen(),
      const JourneyScreen(),
      const ChallengeScreen(),
      const ComingSoonScreen(title: 'Leagues', icon: '🏆', message: 'Weekly leagues are launching soon!'),
      const ComingSoonScreen(title: 'Daily Quests', icon: '📋', message: 'Quest boards with bonus XP and gems are on the way!'),
      const ProfileScreen(),
    ];

    return Scaffold(
      body: screens[_index],
      bottomNavigationBar: GameBottomNav(
        currentIndex: _index,
        onTap: (i) => setState(() => _index = i),
      ),
    );
  }
}
