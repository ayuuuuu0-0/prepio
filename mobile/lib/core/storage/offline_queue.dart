import 'package:hive/hive.dart';

/// PendingAnswer is a queued offline submission.
class PendingAnswer {
  PendingAnswer({
    required this.questionId,
    required this.sessionId,
    required this.answer,
    required this.submittedAt,
  });

  final String questionId;
  final String sessionId;
  final String answer;
  final String submittedAt;

  Map<String, String> toMap() => {
        'question_id': questionId,
        'session_id': sessionId,
        'answer': answer,
        'submitted_at': submittedAt,
      };

  static PendingAnswer fromMap(Map map) => PendingAnswer(
        questionId: map['question_id'] as String,
        sessionId: map['session_id'] as String,
        answer: map['answer'] as String,
        submittedAt: map['submitted_at'] as String,
      );
}

/// OfflineQueue stores pending answer submissions in Hive.
class OfflineQueue {
  OfflineQueue(this._box);

  final Box _box;
  static const boxName = 'pending_answers';

  Future<void> enqueue(PendingAnswer item) async {
    await _box.add(item.toMap());
  }

  List<PendingAnswer> all() {
    return _box.values
        .map((v) => PendingAnswer.fromMap(Map<String, dynamic>.from(v as Map)))
        .toList();
  }

  Future<void> removeAt(int index) async {
    await _box.deleteAt(index);
  }
}
