import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/api/api_client.dart';
import '../../core/config/constants.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/game_button.dart';
import '../../widgets/game/speech_bubble.dart';
import '../auth/auth_provider.dart';
import '../auth/profile_provider.dart';

/// OnboardingScreen collects targets, experience, and companion choice.
class OnboardingScreen extends ConsumerStatefulWidget {
  const OnboardingScreen({super.key});

  @override
  ConsumerState<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends ConsumerState<OnboardingScreen> {
  var _step = 1;
  final _targets = <String>{};
  String _experience = '';
  String _companionId = '';
  List<CompanionInfo> _companions = [];
  var _loading = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadCompanions();
  }

  Future<void> _loadCompanions() async {
    try {
      final api = ref.read(apiClientProvider);
      api.token = ref.read(authTokenProvider);
      final companions = await api.getCompanions();
      setState(() => _companions = companions);
    } catch (e) {
      setState(() => _error = e.toString());
    }
  }

  Future<void> _finish() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final api = ref.read(apiClientProvider);
      api.token = ref.read(authTokenProvider);
      await api.completeOnboarding(
        targetCompanies: _targets.toList(),
        experienceLevel: _experience,
        companionId: _companionId,
      );
      ref.invalidate(profileProvider);
    } catch (e) {
      setState(() => _error = e.toString());
    } finally {
      setState(() => _loading = false);
    }
  }

  String get _stepMessage => switch (_step) {
        1 => "Which dream companies are we aiming for? Pick as many as you like!",
        2 => "How much coding have you done? No wrong answers — we'll meet you where you are.",
        _ => "Every adventurer needs a companion! Who's joining your journey?",
      };

  @override
  Widget build(BuildContext context) {
    final selectedCompanion = _companions.where((c) => c.id == _companionId).firstOrNull;

    return Scaffold(
      body: GameBackground(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              _ProgressBar(step: _step),
              const SizedBox(height: 16),
              Row(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  CompanionHero(
                    name: selectedCompanion?.name ?? 'Byte',
                    species: selectedCompanion?.species,
                    size: 80,
                  ),
                  const SizedBox(width: 12),
                  Expanded(child: SpeechBubble(text: _stepMessage)),
                ],
              ),
              if (_error != null) ...[
                const SizedBox(height: 12),
                Text(_error!, style: const TextStyle(color: PrepioColors.orange, fontWeight: FontWeight.w600)),
              ],
              const SizedBox(height: 16),
              Expanded(child: _buildStep()),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildStep() {
    if (_step == 1) {
      return ListView(
        children: [
          Text('🎯 Target Companies', style: GoogleFonts.fredoka(fontSize: 20, fontWeight: FontWeight.w800)),
          const SizedBox(height: 12),
          ...AppConstants.targetCompanies.map((company) {
            final selected = _targets.contains(company);
            return Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: _ChoiceCard(
                label: company,
                selected: selected,
                color: companyRingColors[company] ?? PrepioColors.green,
                onTap: () => setState(() {
                  if (selected) {
                    _targets.remove(company);
                  } else {
                    _targets.add(company);
                  }
                }),
              ),
            );
          }),
          const SizedBox(height: 16),
          GameButton(
            label: 'Continue →',
            onPressed: _targets.isEmpty ? null : () => setState(() => _step = 2),
          ),
        ],
      );
    }

    if (_step == 2) {
      return ListView(
        children: [
          Text('📚 Experience Level', style: GoogleFonts.fredoka(fontSize: 20, fontWeight: FontWeight.w800)),
          const SizedBox(height: 12),
          ...AppConstants.experienceLevels.map((level) {
            return Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: _ChoiceCard(
                label: level.$2,
                selected: _experience == level.$1,
                color: PrepioColors.blue,
                onTap: () => setState(() => _experience = level.$1),
              ),
            );
          }),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: GameButton(
                  label: '← Back',
                  color: PrepioColors.textMuted,
                  shadowColor: const Color(0xFF555555),
                  onPressed: () => setState(() => _step = 1),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: GameButton(
                  label: 'Continue →',
                  onPressed: _experience.isEmpty ? null : () => setState(() => _step = 3),
                ),
              ),
            ],
          ),
        ],
      );
    }

    return ListView(
      children: [
        Text('🐾 Choose Your Companion', style: GoogleFonts.fredoka(fontSize: 20, fontWeight: FontWeight.w800)),
        const SizedBox(height: 12),
        ..._companions.map((c) {
          final visual = companionFor(name: c.name, species: c.species);
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: _CompanionCard(
              name: c.name,
              species: c.species.replaceAll('_', ' '),
              emoji: visual.emoji,
              colors: visual.colors,
              selected: _companionId == c.id,
              onTap: () => setState(() => _companionId = c.id),
            ),
          );
        }),
        const SizedBox(height: 8),
        Row(
          children: [
            Expanded(
              child: GameButton(
                label: '← Back',
                color: PrepioColors.textMuted,
                shadowColor: const Color(0xFF555555),
                onPressed: () => setState(() => _step = 2),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: GameButton(
                label: _loading ? 'Starting...' : 'Start Journey! 🚀',
                color: PrepioColors.gold,
                shadowColor: const Color(0xFFE5B000),
                onPressed: _companionId.isEmpty || _loading ? null : _finish,
                loading: _loading,
              ),
            ),
          ],
        ),
      ],
    );
  }
}

class _ProgressBar extends StatelessWidget {
  const _ProgressBar({required this.step});
  final int step;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Step $step of 3', style: GoogleFonts.fredoka(fontSize: 14, fontWeight: FontWeight.w700, color: PrepioColors.textMuted)),
        const SizedBox(height: 8),
        Row(
          children: List.generate(3, (i) {
            final active = i < step;
            return Expanded(
              child: Container(
                height: 8,
                margin: EdgeInsets.only(right: i < 2 ? 6 : 0),
                decoration: BoxDecoration(
                  color: active ? PrepioColors.green : Colors.white,
                  borderRadius: BorderRadius.circular(4),
                  border: Border.all(color: active ? PrepioColors.greenDark : PrepioColors.textMuted, width: 2),
                ),
              ),
            );
          }),
        ),
      ],
    );
  }
}

class _ChoiceCard extends StatelessWidget {
  const _ChoiceCard({required this.label, required this.selected, required this.color, required this.onTap});
  final String label;
  final bool selected;
  final Color color;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        decoration: BoxDecoration(
          color: selected ? color.withValues(alpha: 0.15) : Colors.white,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: selected ? color : const Color(0xFFE5E5E5), width: selected ? 3 : 2),
          boxShadow: selected ? [BoxShadow(color: color.withValues(alpha: 0.3), blurRadius: 8, offset: const Offset(0, 4))] : null,
        ),
        child: Row(
          children: [
            Expanded(child: Text(label, style: GoogleFonts.nunito(fontSize: 16, fontWeight: FontWeight.w700))),
            if (selected) Icon(Icons.check_circle, color: color),
          ],
        ),
      ),
    );
  }
}

class _CompanionCard extends StatelessWidget {
  const _CompanionCard({
    required this.name,
    required this.species,
    required this.emoji,
    required this.colors,
    required this.selected,
    required this.onTap,
  });
  final String name;
  final String species;
  final String emoji;
  final List<Color> colors;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(20),
          border: Border.all(color: selected ? PrepioColors.green : const Color(0xFFE5E5E5), width: selected ? 3 : 2),
          boxShadow: selected ? [const BoxShadow(color: Color(0x4058CC02), blurRadius: 12, offset: Offset(0, 4))] : null,
        ),
        child: Row(
          children: [
            Container(
              width: 56,
              height: 56,
              decoration: BoxDecoration(
                gradient: LinearGradient(colors: colors),
                shape: BoxShape.circle,
              ),
              alignment: Alignment.center,
              child: Text(emoji, style: const TextStyle(fontSize: 28)),
            ),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(name, style: GoogleFonts.fredoka(fontSize: 18, fontWeight: FontWeight.w800)),
                  Text(species, style: GoogleFonts.nunito(color: PrepioColors.textMuted, fontWeight: FontWeight.w600)),
                ],
              ),
            ),
            if (selected) const Icon(Icons.favorite, color: PrepioColors.pink),
          ],
        ),
      ),
    );
  }
}
