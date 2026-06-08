import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/hud_bar.dart';
import '../../widgets/game/quest_card.dart';
import '../../widgets/game/readiness_ring.dart';
import '../../widgets/game/speech_bubble.dart';
import 'dashboard_provider.dart';

/// DashboardScreen is the primary home — emotion-first, not analytics.
class DashboardScreen extends ConsumerWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final dashboard = ref.watch(dashboardProvider);

    return GameBackground(
      child: dashboard.when(
        loading: () => const Center(child: CircularProgressIndicator(color: PrepioColors.accent)),
        error: (e, _) => Center(child: Text('$e', style: const TextStyle(color: PrepioColors.danger))),
        data: (home) => ListView(
          padding: const EdgeInsets.fromLTRB(20, 16, 20, 24),
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                CompanionHero(name: home.companion?.name, species: home.companion?.species, size: 80),
                const SizedBox(width: 12),
                Expanded(
                  child: Padding(
                    padding: const EdgeInsets.only(top: 8),
                    child: SpeechBubble(
                      text: home.companionMessage,
                      speakerName: home.companion?.name ?? 'Byte',
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            HudBar(home: home),
            const SizedBox(height: 16),
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: PrepioColors.surface,
                borderRadius: BorderRadius.circular(16),
                border: const Border(
                  left: BorderSide(color: PrepioColors.xp, width: 3),
                  top: BorderSide(color: PrepioColors.border),
                  right: BorderSide(color: PrepioColors.border),
                  bottom: BorderSide(color: PrepioColors.border),
                ),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('🧭 Career Readiness', style: GoogleFonts.plusJakartaSans(fontSize: 16, fontWeight: FontWeight.w700, color: PrepioColors.textPrimary)),
                  const SizedBox(height: 16),
                  Wrap(
                    spacing: 16,
                    runSpacing: 16,
                    alignment: WrapAlignment.center,
                    children: home.readiness
                        .map((r) => ReadinessRing(
                              company: r.company,
                              score: r.score,
                              color: companyRingColors[r.company] ?? PrepioColors.accent,
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
                color: PrepioColors.surface,
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: PrepioColors.border),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('🏆 ${home.league.label}', style: GoogleFonts.plusJakartaSans(color: PrepioColors.textPrimary, fontWeight: FontWeight.w700)),
                  if (home.league.available)
                    Text('Rank #${home.league.rank}', style: GoogleFonts.plusJakartaSans(color: PrepioColors.textPrimary, fontSize: 24, fontWeight: FontWeight.w800))
                  else
                    Text(
                      'Your league rank is calculating — check back soon!',
                      style: GoogleFonts.nunito(color: PrepioColors.textMuted, fontWeight: FontWeight.w600),
                    ),
                ],
              ),
            ),
            const SizedBox(height: 20),
            Text('DAILY QUESTS', style: GoogleFonts.jetBrainsMono(fontSize: 11, fontWeight: FontWeight.w700, color: PrepioColors.textDim, letterSpacing: 1.2)),
            const SizedBox(height: 8),
            ...home.dailyQuests.map((q) => QuestCard(
                  title: q.title,
                  icon: _questIcon(q.id),
                  progress: q.progress,
                  target: q.target,
                  completed: q.completed,
                  rewardXp: q.rewardXp,
                  rewardGems: q.rewardGems,
                  comingSoon: q.comingSoon,
                )),
          ],
        ),
      ),
    );
  }

  String _questIcon(String id) {
    if (id.contains('streak')) return '🔥';
    if (id.contains('score')) return '🎯';
    return '⚡';
  }
}
