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
import '../question/question_provider.dart';

/// ChallengeScreen presents today's questions with inline dark result card.
class ChallengeScreen extends ConsumerStatefulWidget {
  const ChallengeScreen({super.key});

  @override
  ConsumerState<ChallengeScreen> createState() => _ChallengeScreenState();
}

class _ChallengeScreenState extends ConsumerState<ChallengeScreen> {
  final _answer = TextEditingController();
  var _submitting = false;
  var _currentIndex = 0;
  final _answeredIds = <String>{};
  SubmitResult? _result;

  @override
  void dispose() {
    _answer.dispose();
    super.dispose();
  }

  int get _trimmedLength => _answer.text.trim().length;
  bool get _canSubmit => _trimmedLength >= AppConstants.minAnswerLength;

  Future<void> _loadHistory(String sessionId) async {
    final api = ref.read(apiClientProvider);
    api.token = ref.read(authTokenProvider);
    final history = await api.getQuestionHistory(sessionId);
    setState(() {
      _answeredIds
        ..clear()
        ..addAll(history.map((h) => h.questionId));
      _currentIndex = 0;
    });
  }

  Future<void> _submit(String questionId, String sessionId) async {
    if (!_canSubmit) return;
    setState(() => _submitting = true);
    try {
      final api = ref.read(apiClientProvider);
      api.token = ref.read(authTokenProvider);
      final result = await api.submitAnswer(questionId, sessionId, _answer.text.trim());
      setState(() {
        _result = result;
        _answeredIds.add(questionId);
        _answer.clear();
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
      }
    } finally {
      setState(() => _submitting = false);
    }
  }

  void _nextQuestion(List<QuestionItem> questions) {
    setState(() {
      _result = null;
      for (var i = _currentIndex + 1; i < questions.length; i++) {
        if (!_answeredIds.contains(questions[i].id)) {
          _currentIndex = i;
          return;
        }
      }
    });
  }

  bool _hasNext(List<QuestionItem> questions) {
    for (var i = _currentIndex + 1; i < questions.length; i++) {
      if (!_answeredIds.contains(questions[i].id)) return true;
    }
    return false;
  }

  String _questionHint(QuestionItem question) {
    const hints = {
      'dsa': 'Think about time and space complexity before you write.',
      'system_design': 'Start with requirements, then components, then trade-offs.',
      'lld': 'Focus on class design, responsibilities, and relationships.',
      'behavioral': 'Use the STAR format: Situation, Task, Action, Result.',
      'aptitude': 'Break the problem into smaller steps first.',
      'fundamentals': 'Cover the core concept clearly before diving into detail.',
    };
    return hints[question.roundType] ?? 'Take your time. Think it through.';
  }

  @override
  Widget build(BuildContext context) {
    final paperAsync = ref.watch(dailyPaperProvider);
    final profile = ref.watch(profileProvider);

    return paperAsync.when(
      loading: () => const Scaffold(
        backgroundColor: PrepioColors.bg,
        body: Center(child: CircularProgressIndicator(color: PrepioColors.accent)),
      ),
      error: (e, _) => Scaffold(backgroundColor: PrepioColors.bg, body: Center(child: Text('$e'))),
      data: (paper) {
        if (paper.questions.isEmpty) {
          return const Scaffold(backgroundColor: PrepioColors.bg, body: Center(child: Text('No challenge today')));
        }

        if (_answeredIds.isEmpty && paper.sessionId.isNotEmpty) {
          _loadHistory(paper.sessionId);
        }

        while (_currentIndex < paper.questions.length && _answeredIds.contains(paper.questions[_currentIndex].id)) {
          _currentIndex++;
        }
        if (_currentIndex >= paper.questions.length) _currentIndex = paper.questions.length - 1;

        final question = paper.questions[_currentIndex];
        final companionName = profile.valueOrNull?.companion?.name ?? 'Byte';
        final companionSpecies = profile.valueOrNull?.companion?.species;
        final roundColor = roundTypeColors[question.roundType] ?? PrepioColors.accent;

        return Scaffold(
          backgroundColor: PrepioColors.bg,
          body: GameBackground(
            variant: GameBgVariant.challenge,
            child: ListView(
              padding: const EdgeInsets.all(20),
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text('← Home', style: GoogleFonts.jetBrainsMono(fontWeight: FontWeight.w600, color: PrepioColors.textDim)),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                      decoration: BoxDecoration(
                        color: PrepioColors.surface,
                        borderRadius: BorderRadius.circular(8),
                        border: Border.all(color: PrepioColors.border),
                      ),
                      child: Text(
                        'Q ${_currentIndex + 1}/${paper.questions.length}',
                        style: GoogleFonts.jetBrainsMono(fontWeight: FontWeight.w700, color: PrepioColors.textMuted),
                      ),
                    ),
                    CompanionHero(
                      name: companionName,
                      species: companionSpecies,
                      size: 48,
                      reaction: _result != null
                          ? (_result!.correct ? CompanionReaction.correct : CompanionReaction.wrong)
                          : CompanionReaction.idle,
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                SpeechBubble(speakerName: companionName, text: _questionHint(question)),
                const SizedBox(height: 16),
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    color: PrepioColors.surface,
                    borderRadius: BorderRadius.circular(16),
                    border: Border(
                      left: BorderSide(color: roundColor, width: 3),
                      top: const BorderSide(color: PrepioColors.border),
                      right: const BorderSide(color: PrepioColors.border),
                      bottom: const BorderSide(color: PrepioColors.border),
                    ),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Wrap(
                        spacing: 8,
                        children: [
                          _tag(question.roundType.replaceAll('_', ' '), roundColor),
                          _tag(question.difficulty, _difficultyColor(question.difficulty)),
                        ],
                      ),
                      const SizedBox(height: 12),
                      Text(question.body, style: GoogleFonts.nunito(fontSize: 16, fontWeight: FontWeight.w500, height: 1.5, color: PrepioColors.textBody)),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
                if (_result != null)
                  _ResultCard(
                    result: _result!,
                    hasNext: _hasNext(paper.questions),
                    onNext: () => _nextQuestion(paper.questions),
                  )
                else ...[
                  TextField(
                    controller: _answer,
                    maxLines: 8,
                    onChanged: (_) => setState(() {}),
                    style: const TextStyle(color: PrepioColors.textPrimary),
                    decoration: const InputDecoration(hintText: 'Explain your approach...'),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    _canSubmit
                        ? 'Ready to submit'
                        : '$_trimmedLength/${AppConstants.minAnswerLength} chars — add more detail',
                    style: GoogleFonts.jetBrainsMono(
                      fontSize: 12,
                      fontWeight: FontWeight.w600,
                      color: _canSubmit ? PrepioColors.success : PrepioColors.textDim,
                    ),
                  ),
                  const SizedBox(height: 16),
                  GameButton(
                    label: _submitting ? 'Evaluating...' : 'Submit',
                    onPressed: (_submitting || !_canSubmit) ? null : () => _submit(question.id, paper.sessionId),
                    loading: _submitting,
                  ),
                ],
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _tag(String label, Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(label.toUpperCase(), style: GoogleFonts.jetBrainsMono(fontSize: 10, fontWeight: FontWeight.w700, color: color)),
    );
  }

  Color _difficultyColor(String difficulty) {
    if (difficulty == 'hard') return PrepioColors.danger;
    if (difficulty == 'medium') return PrepioColors.warning;
    return PrepioColors.success;
  }
}

class _ResultCard extends StatelessWidget {
  const _ResultCard({required this.result, required this.hasNext, required this.onNext});

  final SubmitResult result;
  final bool hasNext;
  final VoidCallback onNext;

  @override
  Widget build(BuildContext context) {
    final accent = result.correct ? PrepioColors.success : PrepioColors.warning;

    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: accent.withValues(alpha: 0.08),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: accent.withValues(alpha: 0.3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            result.correct ? 'Solid.' : 'Not quite.',
            style: GoogleFonts.plusJakartaSans(fontSize: 20, fontWeight: FontWeight.w800, color: accent),
          ),
          const SizedBox(height: 4),
          Text(result.feedback, style: GoogleFonts.nunito(color: PrepioColors.textMuted)),
          Text('Score: ${result.score}%', style: GoogleFonts.jetBrainsMono(fontSize: 12, color: PrepioColors.textDim)),
          if (result.strengths.isNotEmpty) ...[
            const SizedBox(height: 12),
            Text('COVERED', style: GoogleFonts.jetBrainsMono(fontSize: 10, color: PrepioColors.textDim)),
            ...result.strengths.map((s) => Text('✓ $s', style: const TextStyle(color: PrepioColors.textBody))),
          ],
          if (result.gaps.isNotEmpty) ...[
            const SizedBox(height: 8),
            Text('MISSING', style: GoogleFonts.jetBrainsMono(fontSize: 10, color: PrepioColors.textDim)),
            ...result.gaps.map((g) => Text('→ $g', style: const TextStyle(color: PrepioColors.textBody))),
          ],
          if (result.correct) ...[
            const SizedBox(height: 12),
            Row(
              children: [
                _rewardChip('+${result.xpAwarded} XP', PrepioColors.xp),
                const SizedBox(width: 8),
                _rewardChip('+${result.gemsAwarded} 💎', PrepioColors.gems),
              ],
            ),
          ],
          if (result.readinessDelta > 0)
            Padding(
              padding: const EdgeInsets.only(top: 8),
              child: Text('Readiness +${result.readinessDelta}%', style: GoogleFonts.jetBrainsMono(fontSize: 12, color: PrepioColors.xp)),
            ),
          const SizedBox(height: 16),
          GameButton(label: hasNext ? 'Next Question →' : 'Back to dashboard', onPressed: hasNext ? onNext : null),
        ],
      ),
    );
  }

  Widget _rewardChip(String text, Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withValues(alpha: 0.2)),
      ),
      child: Text(text, style: GoogleFonts.jetBrainsMono(fontWeight: FontWeight.w700, color: color)),
    );
  }
}
