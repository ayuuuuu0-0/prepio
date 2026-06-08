import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';
import '../../core/config/constants.dart';
import '../../core/theme/design_tokens.dart';
import '../../widgets/game/companion_hero.dart';
import '../../widgets/game/game_background.dart';
import '../../widgets/game/game_button.dart';
import '../../widgets/game/speech_bubble.dart';
import '../auth/auth_provider.dart';
import '../question/question_provider.dart';

/// ChallengeScreen presents today's question in focus mode with celebration results.
class ChallengeScreen extends ConsumerStatefulWidget {
  const ChallengeScreen({super.key});

  @override
  ConsumerState<ChallengeScreen> createState() => _ChallengeScreenState();
}

class _ChallengeScreenState extends ConsumerState<ChallengeScreen> {
  final _answer = TextEditingController();
  var _submitting = false;
  String? _feedback;
  bool? _correct;

  @override
  void dispose() {
    _answer.dispose();
    super.dispose();
  }

  int get _trimmedLength => _answer.text.trim().length;

  bool get _canSubmit => _trimmedLength >= AppConstants.minAnswerLength;

  Future<void> _submit(String questionId, String sessionId) async {
    if (!_canSubmit) {
      setState(() {
        _feedback = 'Write at least ${AppConstants.minAnswerLength} characters before submitting.';
        _correct = null;
      });
      return;
    }
    setState(() {
      _submitting = true;
      _feedback = null;
      _correct = null;
    });
    try {
      final api = ref.read(apiClientProvider);
      api.token = ref.read(authTokenProvider);
      final result = await api.submitAnswer(questionId, sessionId, _answer.text.trim());
      setState(() {
        _feedback = result.feedback;
        _correct = result.correct;
      });
      _answer.clear();
    } catch (e) {
      final message = e.toString();
      setState(() {
        _correct = null;
        _feedback = message.contains('already submitted')
            ? 'You already answered this question today. Come back tomorrow!'
            : message;
      });
    } finally {
      setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final paperAsync = ref.watch(dailyPaperProvider);

    return Scaffold(
      body: GameBackground(
        variant: GameBgVariant.challenge,
        child: paperAsync.when(
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (e, _) => Center(child: Text('$e')),
          data: (paper) {
            if (paper.questions.isEmpty) return const Center(child: Text('No challenge today'));
            final question = paper.questions.first;

            return ListView(
              padding: const EdgeInsets.all(20),
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    IconButton(icon: const Icon(Icons.arrow_back, color: PrepioColors.blue), onPressed: () => Navigator.pop(context)),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(20)),
                      child: Text('Q 1/5', style: GoogleFonts.fredoka(fontWeight: FontWeight.w700)),
                    ),
                    const CompanionHero(size: 48),
                  ],
                ),
                if (_correct != null) ...[
                  const SizedBox(height: 16),
                  Container(
                    padding: const EdgeInsets.all(24),
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        colors: _correct!
                            ? [PrepioColors.green, const Color(0xFF7CB342)]
                            : [const Color(0xFFFFB84D), PrepioColors.orange],
                      ),
                      borderRadius: BorderRadius.circular(24),
                    ),
                    child: Column(
                      children: [
                        Text(_correct! ? '🎉' : '💪', style: const TextStyle(fontSize: 48)),
                        Text(_correct! ? 'Amazing!' : 'Almost there!', style: GoogleFonts.fredoka(fontSize: 24, fontWeight: FontWeight.w800, color: Colors.white)),
                        const SizedBox(height: 8),
                        Text(_feedback ?? '', textAlign: TextAlign.center, style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w600)),
                        if (_correct!) ...[
                          const SizedBox(height: 12),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              _rewardChip('⚡ +XP'),
                              const SizedBox(width: 12),
                              _rewardChip('💎 +Gems'),
                            ],
                          ),
                        ],
                      ],
                    ),
                  ),
                ],
                if (_correct == null) ...[
                  if (_feedback != null) ...[
                    const SizedBox(height: 16),
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: const Color(0xFFFFF3E0),
                        borderRadius: BorderRadius.circular(16),
                        border: Border.all(color: PrepioColors.orange, width: 2),
                      ),
                      child: Text(_feedback!, style: GoogleFonts.nunito(fontWeight: FontWeight.w700, color: PrepioColors.orange)),
                    ),
                  ],
                  const SizedBox(height: 16),
                  const SpeechBubble(text: "You've got this! Take your time and think it through."),
                  const SizedBox(height: 16),
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
                        Chip(label: Text(question.difficulty.toUpperCase(), style: GoogleFonts.fredoka(fontWeight: FontWeight.w700))),
                        const SizedBox(height: 12),
                        Text(question.body, style: GoogleFonts.nunito(fontSize: 16, fontWeight: FontWeight.w600, height: 1.5)),
                        const SizedBox(height: 16),
                        TextField(
                          controller: _answer,
                          maxLines: 5,
                          onChanged: (_) => setState(() {}),
                          decoration: const InputDecoration(hintText: 'Write your answer...'),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          _canSubmit
                              ? 'Ready to submit!'
                              : '$_trimmedLength/${AppConstants.minAnswerLength} characters — keep writing...',
                          style: GoogleFonts.nunito(
                            fontSize: 12,
                            fontWeight: FontWeight.w700,
                            color: _canSubmit ? PrepioColors.green : PrepioColors.textMuted,
                          ),
                        ),
                        const SizedBox(height: 16),
                        GameButton(
                          label: _submitting ? 'Checking...' : 'Submit Answer!',
                          onPressed: (_submitting || !_canSubmit) ? null : () => _submit(question.id, paper.sessionId),
                          loading: _submitting,
                        ),
                      ],
                    ),
                  ),
                ],
              ],
            );
          },
        ),
      ),
    );
  }

  Widget _rewardChip(String text) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(color: Colors.white24, borderRadius: BorderRadius.circular(20)),
      child: Text(text, style: GoogleFonts.fredoka(fontWeight: FontWeight.w700, color: Colors.white)),
    );
  }
}
