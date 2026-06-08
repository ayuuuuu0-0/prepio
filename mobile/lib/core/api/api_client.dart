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

    if (res.statusCode >= 400) {
      throw ApiException(_parseErrorMessage(res));
    }

    final decoded = _decodeBody(res.body);
    final data = decoded['data'];
    if (data is! Map<String, dynamic>) {
      throw ApiException('unexpected response from $path');
    }
    return data;
  }

  Future<List<dynamic>> _requestList(String path) async {
    final headers = <String, String>{'Content-Type': 'application/json'};
    if (token != null && token!.isNotEmpty) {
      headers['Authorization'] = 'Bearer $token';
    }

    final res = await http.get(Uri.parse('$baseUrl$path'), headers: headers);
    if (res.statusCode >= 400) {
      throw ApiException(_parseErrorMessage(res));
    }

    final decoded = _decodeBody(res.body);
    final data = decoded['data'];
    if (data is! List<dynamic>) {
      throw ApiException('unexpected response from $path');
    }
    return data;
  }

  Map<String, dynamic> _decodeBody(String body) {
    if (body.isEmpty) {
      throw ApiException('empty response from server');
    }
    try {
      final decoded = jsonDecode(body);
      if (decoded is! Map<String, dynamic>) {
        throw ApiException('invalid response from server');
      }
      return decoded;
    } catch (_) {
      throw ApiException('server unavailable — restart backend with: make dev');
    }
  }

  String _parseErrorMessage(http.Response res) {
    try {
      final decoded = jsonDecode(res.body);
      if (decoded is Map<String, dynamic>) {
        final err = decoded['error'];
        if (err is Map<String, dynamic> && err['message'] is String) {
          return err['message'] as String;
        }
      }
    } catch (_) {}
    return 'request failed (${res.statusCode})';
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

  Future<ProfileInfo> getProfile() async {
    final data = await _request('GET', '/api/v1/users/profile');
    return ProfileInfo.fromJson(data);
  }

  Future<List<CompanionInfo>> getCompanions() async {
    final list = await _requestList('/api/v1/companions');
    return list.map((e) => CompanionInfo.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<ProfileInfo> completeOnboarding({
    required List<String> targetCompanies,
    required String experienceLevel,
    required String companionId,
  }) async {
    final data = await _request('POST', '/api/v1/users/onboarding', body: {
      'target_companies': targetCompanies,
      'experience_level': experienceLevel,
      'companion_id': companionId,
    });
    return ProfileInfo.fromJson(data);
  }

  Future<DashboardHome> getDashboardHome() async {
    final data = await _request('GET', '/api/v1/dashboard/home');
    return DashboardHome.fromJson(data);
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

class CompanionInfo {
  CompanionInfo({required this.id, required this.name, required this.species});
  final String id;
  final String name;
  final String species;

  factory CompanionInfo.fromJson(Map<String, dynamic> json) {
    return CompanionInfo(
      id: json['id'] as String,
      name: json['name'] as String,
      species: json['species'] as String,
    );
  }
}

class ProfileInfo {
  ProfileInfo({required this.onboardingCompleted, required this.targetCompanies, this.companion});
  final bool onboardingCompleted;
  final List<String> targetCompanies;
  final CompanionInfo? companion;

  factory ProfileInfo.fromJson(Map<String, dynamic> json) {
    CompanionInfo? companion;
    if (json['companion'] != null) {
      companion = CompanionInfo.fromJson(json['companion'] as Map<String, dynamic>);
    }
    final targets = json['target_companies'];
    return ProfileInfo(
      onboardingCompleted: json['onboarding_completed'] as bool? ?? false,
      targetCompanies: targets is List<dynamic> ? targets.cast<String>() : <String>[],
      companion: companion,
    );
  }
}

/// LeagueInfo summarizes the user's league tier and rank.
class LeagueInfo {
  LeagueInfo({required this.tier, required this.rank, required this.label});
  final String tier;
  final int rank;
  final String label;

  factory LeagueInfo.fromJson(Map<String, dynamic> json) {
    return LeagueInfo(
      tier: json['tier'] as String,
      rank: json['rank'] as int,
      label: json['label'] as String,
    );
  }
}

class DashboardHome {
  DashboardHome({
    required this.streak,
    required this.progress,
    required this.companionMessage,
    required this.readiness,
    required this.league,
    required this.dailyQuests,
    required this.onboardingNeeded,
    this.companion,
  });

  final StreakInfo streak;
  final ProgressInfo progress;
  final String companionMessage;
  final List<ReadinessInfo> readiness;
  final LeagueInfo league;
  final List<DailyQuestInfo> dailyQuests;
  final bool onboardingNeeded;
  final CompanionInfo? companion;

  factory DashboardHome.fromJson(Map<String, dynamic> json) {
    CompanionInfo? companion;
    if (json['companion'] != null) {
      companion = CompanionInfo.fromJson(json['companion'] as Map<String, dynamic>);
    }
    return DashboardHome(
      streak: StreakInfo.fromJson(json['streak'] as Map<String, dynamic>),
      progress: ProgressInfo.fromJson(json['progress'] as Map<String, dynamic>),
      companionMessage: json['companion_message'] as String,
      readiness: (json['readiness'] as List<dynamic>)
          .map((e) => ReadinessInfo.fromJson(e as Map<String, dynamic>))
          .toList(),
      league: LeagueInfo.fromJson(json['league'] as Map<String, dynamic>),
      dailyQuests: (json['daily_quests'] as List<dynamic>)
          .map((e) => DailyQuestInfo.fromJson(e as Map<String, dynamic>))
          .toList(),
      onboardingNeeded: json['onboarding_needed'] as bool,
      companion: companion,
    );
  }
}

class ReadinessInfo {
  ReadinessInfo({required this.company, required this.score});
  final String company;
  final int score;

  factory ReadinessInfo.fromJson(Map<String, dynamic> json) {
    return ReadinessInfo(company: json['company'] as String, score: json['score'] as int);
  }
}

class DailyQuestInfo {
  DailyQuestInfo({required this.title, required this.progress, required this.target, required this.completed});
  final String title;
  final int progress;
  final int target;
  final bool completed;

  factory DailyQuestInfo.fromJson(Map<String, dynamic> json) {
    return DailyQuestInfo(
      title: json['title'] as String,
      progress: json['progress'] as int,
      target: json['target'] as int,
      completed: json['completed'] as bool,
    );
  }
}
