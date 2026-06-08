import 'dart:convert';
import 'package:http/http.dart' as http;

/// ApiClient talks to the Prepio gateway REST API.
class ApiClient {
  ApiClient({this.baseUrl = 'http://localhost:8080', this.token});

  final String baseUrl;
  String? token;

  Future<Map<String, dynamic>> _request(
    String method,
    String path, {
    Map<String, dynamic>? body,
  }) async {
    final headers = <String, String>{'Content-Type': 'application/json'};
    if (token != null && token!.isNotEmpty) {
      headers['Authorization'] = 'Bearer $token';
    }

    final uri = Uri.parse('$baseUrl$path');
    late http.Response res;

    switch (method) {
      case 'GET':
        res = await http.get(uri, headers: headers);
      case 'POST':
        res = await http.post(uri, headers: headers, body: jsonEncode(body));
      default:
        throw UnsupportedError(method);
    }

    final decoded = jsonDecode(res.body) as Map<String, dynamic>;
    if (res.statusCode >= 400) {
      final err = decoded['error'] as Map<String, dynamic>?;
      throw ApiException(err?['message'] as String? ?? 'request failed');
    }
    return decoded['data'] as Map<String, dynamic>;
  }

  Future<AuthResult> register(String email, String username, String password) async {
    final data = await _request('POST', '/api/v1/auth/register', body: {
      'email': email,
      'username': username,
      'password': password,
    });
    return AuthResult.fromJson(data);
  }

  Future<AuthResult> login(String email, String password) async {
    final data = await _request('POST', '/api/v1/auth/login', body: {
      'email': email,
      'password': password,
    });
    return AuthResult.fromJson(data);
  }

  Future<DailyPaper> getDailyPaper() async {
    final data = await _request('GET', '/api/v1/questions/daily');
    return DailyPaper.fromJson(data);
  }

  Future<SubmitResult> submitAnswer(String questionId, String sessionId, String answer) async {
    final data = await _request('POST', '/api/v1/questions/$questionId/submit', body: {
      'session_id': sessionId,
      'answer': answer,
      'time_spent_seconds': 60,
    });
    return SubmitResult.fromJson(data);
  }

  Future<StreakInfo> getStreak() async {
    final data = await _request('GET', '/api/v1/streaks/me');
    return StreakInfo.fromJson(data);
  }

  Future<ProgressInfo> getProgress() async {
    final data = await _request('GET', '/api/v1/progress/me');
    return ProgressInfo.fromJson(data);
  }
}

class ApiException implements Exception {
  ApiException(this.message);
  final String message;

  @override
  String toString() => message;
}

class AuthResult {
  AuthResult({required this.accessToken, required this.username});
  final String accessToken;
  final String username;

  factory AuthResult.fromJson(Map<String, dynamic> json) {
    final user = json['user'] as Map<String, dynamic>;
    return AuthResult(
      accessToken: json['access_token'] as String,
      username: user['username'] as String,
    );
  }
}

class DailyPaper {
  DailyPaper({required this.sessionId, required this.questions});
  final String sessionId;
  final List<QuestionItem> questions;

  factory DailyPaper.fromJson(Map<String, dynamic> json) {
    final qs = (json['questions'] as List<dynamic>)
        .map((q) => QuestionItem.fromJson(q as Map<String, dynamic>))
        .toList();
    return DailyPaper(sessionId: json['session_id'] as String, questions: qs);
  }
}

class QuestionItem {
  QuestionItem({required this.id, required this.body, required this.difficulty});
  final String id;
  final String body;
  final String difficulty;

  factory QuestionItem.fromJson(Map<String, dynamic> json) {
    return QuestionItem(
      id: json['id'] as String,
      body: json['body'] as String,
      difficulty: json['difficulty'] as String,
    );
  }
}

class SubmitResult {
  SubmitResult({required this.correct, required this.feedback});
  final bool correct;
  final String feedback;

  factory SubmitResult.fromJson(Map<String, dynamic> json) {
    return SubmitResult(
      correct: json['correct'] as bool,
      feedback: json['feedback'] as String,
    );
  }
}

class StreakInfo {
  StreakInfo({required this.currentStreak, required this.streakActiveToday});
  final int currentStreak;
  final bool streakActiveToday;

  factory StreakInfo.fromJson(Map<String, dynamic> json) {
    return StreakInfo(
      currentStreak: json['current_streak'] as int,
      streakActiveToday: json['streak_active_today'] as bool,
    );
  }
}

class ProgressInfo {
  ProgressInfo({required this.totalXp, required this.currentLevel, required this.gemBalance});
  final int totalXp;
  final int currentLevel;
  final int gemBalance;

  factory ProgressInfo.fromJson(Map<String, dynamic> json) {
    return ProgressInfo(
      totalXp: json['total_xp'] as int,
      currentLevel: json['current_level'] as int,
      gemBalance: json['gem_balance'] as int,
    );
  }
}
