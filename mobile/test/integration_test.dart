import 'package:flutter_test/flutter_test.dart';
import 'package:prepio/core/api/api_client.dart';

/// Live API integration test — requires backend at localhost:8080.
void main() {
  test('mobile api client full loop', () async {
    final client = ApiClient();
    final suffix = DateTime.now().millisecondsSinceEpoch;

    final auth = await client.register(
      'mobile-$suffix@test.com',
      'mobile_$suffix',
      'password123',
    );
    expect(auth.accessToken.isNotEmpty, true);

    client.token = auth.accessToken;
    final daily = await client.getDailyPaper();
    expect(daily.questions.isNotEmpty, true);

    final submit = await client.submitAnswer(
      daily.questions.first.id,
      daily.sessionId,
      'hash map approach with O(n) time and O(n) space complexity',
    );
    expect(submit.correct, true);

    final streak = await client.getStreak();
    expect(streak.currentStreak, greaterThanOrEqualTo(1));

    final progress = await client.getProgress();
    expect(progress.totalXp, greaterThan(0));
  }, skip: false);
}
