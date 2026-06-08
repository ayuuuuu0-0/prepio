import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../auth/auth_provider.dart';
import '../progress/progress_provider.dart';
import '../question/question_provider.dart';
import '../streak/streak_provider.dart';
import '../../core/api/api_client.dart';

/// HomeScreen shows daily question, streak, and progress.
class HomeScreen extends ConsumerStatefulWidget {
  const HomeScreen({super.key});

  @override
  ConsumerState<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends ConsumerState<HomeScreen> {
  final _answer = TextEditingController();
  var _submitting = false;
  String? _feedback;

  @override
  void dispose() {
    _answer.dispose();
    super.dispose();
  }

  Future<void> _submit(DailyPaper paper, QuestionItem question) async {
    if (_answer.text.length < 10) return;
    setState(() => _submitting = true);
    try {
      final api = ref.read(apiClientProvider);
      api.token = ref.read(authTokenProvider);
      final result = await api.submitAnswer(question.id, paper.sessionId, _answer.text);
      setState(() => _feedback = result.feedback);
      _answer.clear();
      ref.invalidate(dailyPaperProvider);
      ref.invalidate(streakProvider);
      ref.invalidate(progressProvider);
    } catch (e) {
      setState(() => _feedback = e.toString());
    } finally {
      setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final paperAsync = ref.watch(dailyPaperProvider);
    final streakAsync = ref.watch(streakProvider);
    final progressAsync = ref.watch(progressProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Prepio'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () => ref.read(authTokenProvider.notifier).state = null,
          ),
        ],
      ),
      body: paperAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('$e')),
        data: (paper) {
          if (paper.questions.isEmpty) {
            return const Center(child: Text('No questions today'));
          }
          final question = paper.questions.first;
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              streakAsync.when(
                data: (s) => Card(
                  child: ListTile(
                    title: const Text('Streak'),
                    subtitle: Text(
                      s.streakActiveToday ? 'Active today' : 'Not active yet today',
                    ),
                    trailing: Text('${s.currentStreak} days', style: Theme.of(context).textTheme.titleLarge),
                  ),
                ),
                loading: () => const LinearProgressIndicator(),
                error: (e, _) => Text('$e'),
              ),
              const SizedBox(height: 12),
              progressAsync.when(
                data: (p) => Card(
                  child: ListTile(
                    title: Text('Level ${p.currentLevel}'),
                    subtitle: Text('${p.totalXp} XP · ${p.gemBalance} gems'),
                  ),
                ),
                loading: () => const LinearProgressIndicator(),
                error: (e, _) => Text('$e'),
              ),
              const SizedBox(height: 12),
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text("Today's question", style: Theme.of(context).textTheme.titleMedium),
                      const SizedBox(height: 8),
                      Chip(label: Text(question.difficulty)),
                      const SizedBox(height: 12),
                      Text(question.body),
                      const SizedBox(height: 12),
                      TextField(
                        controller: _answer,
                        maxLines: 4,
                        decoration: const InputDecoration(
                          border: OutlineInputBorder(),
                          hintText: 'Write your answer...',
                        ),
                      ),
                      const SizedBox(height: 12),
                      FilledButton(
                        onPressed: _submitting ? null : () => _submit(paper, question),
                        child: Text(_submitting ? 'Submitting...' : 'Submit answer'),
                      ),
                      if (_feedback != null) ...[
                        const SizedBox(height: 12),
                        Text(_feedback!),
                      ],
                    ],
                  ),
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}
