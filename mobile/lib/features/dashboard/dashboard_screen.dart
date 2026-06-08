import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/game_button.dart';
import '../../widgets/game/quest_card.dart';
import '../../widgets/game/readiness_ring.dart';
import '../../widgets/game/speech_bubble.dart';
import '../auth/auth_provider.dart'; // sessionStoreProvider, authTokenProvider
import '../challenge/challenge_screen.dart';
import '../journey/journey_screen.dart';
import 'dashboard_provider.dart';

/// DashboardScreen is the primary home — emotion-first, not analytics.
class DashboardScreen extends ConsumerWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final dashboard = ref.watch(dashboardProvider);

    return Scaffold(
      body: GameBackground(
        child: dashboard.when(
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (e, _) => Center(child: Text('$e')),
          data: (home) => ListView(
            padding: const EdgeInsets.fromLTRB(20, 16, 20, 100),
            children: [
              Row(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  CompanionHero(name: home.companion?.name, species: home.companion?.species, size: 100),
                  const SizedBox(width: 12),
                  Expanded(child: SpeechBubble(text: home.companionMessage)),
                ],
              ),
              const SizedBox(height: 16),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: [
                  _StatChip(icon: '🔥', label: 'Streak', value: '${home.streak.currentStreak}d', color: PrepioColors.orange),
                  _StatChip(icon: '⚡', label: 'Level', value: '${home.progress.currentLevel}', color: PrepioColors.blue),
                  _StatChip(icon: '💎', label: 'Gems', value: '${home.progress.gemBalance}', color: PrepioColors.purple),
                ],
              ),
              const SizedBox(height: 20),
              Container(
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(24),
                  boxShadow: const [BoxShadow(color: Colors.black12, blurRadius: 8)],
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('🧭 Career Readiness', style: GoogleFonts.fredoka(fontSize: 18, fontWeight: FontWeight.w700)),
                    const SizedBox(height: 16),
                    Wrap(
                      spacing: 16,
                      runSpacing: 16,
                      alignment: WrapAlignment.center,
                      children: home.readiness
                          .map((r) => ReadinessRing(
                                company: r.company,
                                score: r.score,
                                color: companyRingColors[r.company] ?? PrepioColors.green,
                              ))
                          .toList(),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  gradient: const LinearGradient(colors: [Color(0xFF3B82F6), Color(0xFF22D3EE)]),
                  borderRadius: BorderRadius.circular(24),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('🏆 ${home.league.label}', style: GoogleFonts.fredoka(color: Colors.white, fontWeight: FontWeight.w700)),
                    Text('Rank #${home.league.rank}', style: GoogleFonts.fredoka(color: Colors.white, fontSize: 28, fontWeight: FontWeight.w800)),
                  ],
                ),
              ),
              const SizedBox(height: 20),
              Text('⚡ Daily Quests', style: GoogleFonts.fredoka(fontSize: 18, fontWeight: FontWeight.w700)),
              const SizedBox(height: 8),
              ...home.dailyQuests.map((q) => QuestCard(
                    title: q.title,
                    icon: _questIcon(q.title),
                    progress: q.progress,
                    target: q.target,
                    completed: q.completed,
                    rewardXp: 50,
                    rewardGems: 10,
                  )),
              const SizedBox(height: 16),
              GameButton(
                label: 'Continue Journey →',
                onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const ChallengeScreen())),
              ),
              const SizedBox(height: 8),
              GameButton(
                label: 'View Journey Map',
                color: PrepioColors.blue,
                shadowColor: const Color(0xFF1899D6),
                onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const JourneyScreen())),
              ),
              TextButton(
                onPressed: () async {
                  await ref.read(sessionStoreProvider).clearToken();
                  ref.read(authTokenProvider.notifier).state = null;
                },
                child: const Text('Sign out', style: TextStyle(color: PrepioColors.textMuted)),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: _BottomNav(current: 0),
    );
  }

  String _questIcon(String title) {
    if (title.contains('streak')) return '🔥';
    if (title.contains('80%')) return '🎯';
    return '⚡';
  }
}

class _StatChip extends StatelessWidget {
  const _StatChip({required this.icon, required this.label, required this.value, required this.color});
  final String icon;
  final String label;
  final String value;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.15),
        borderRadius: BorderRadius.circular(999),
        border: Border.all(color: color, width: 2),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(icon, style: const TextStyle(fontSize: 18)),
          const SizedBox(width: 6),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(label, style: GoogleFonts.fredoka(fontSize: 10, color: color, fontWeight: FontWeight.w700)),
              Text(value, style: GoogleFonts.fredoka(fontSize: 14, fontWeight: FontWeight.w800)),
            ],
          ),
        ],
      ),
    );
  }
}

class _BottomNav extends StatelessWidget {
  const _BottomNav({required this.current});
  final int current;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        color: Colors.white,
        border: Border(top: BorderSide(color: Color(0xFFE5E5E5), width: 3)),
      ),
      child: SafeArea(
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            _NavItem(icon: '🏠', label: 'Home', active: current == 0),
            _NavItem(icon: '🗺️', label: 'Journey', active: current == 1),
            _NavItem(icon: '⚡', label: 'Play', active: current == 2),
          ],
        ),
      ),
    );
  }
}

class _NavItem extends StatelessWidget {
  const _NavItem({required this.icon, required this.label, required this.active});
  final String icon;
  final String label;
  final bool active;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(icon, style: TextStyle(fontSize: 22, color: active ? PrepioColors.green : PrepioColors.textMuted)),
          Text(label, style: GoogleFonts.fredoka(fontSize: 11, fontWeight: FontWeight.w700, color: active ? PrepioColors.green : PrepioColors.textMuted)),
        ],
      ),
    );
  }
}
