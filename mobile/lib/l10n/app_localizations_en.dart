// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'Prepio';

  @override
  String get signIn => 'Sign in';

  @override
  String get register => 'Register';

  @override
  String get email => 'Email';

  @override
  String get password => 'Password';

  @override
  String get username => 'Username';

  @override
  String get dailyQuestion => 'Today\'s question';

  @override
  String get submitAnswer => 'Submit answer';

  @override
  String get streak => 'Streak';

  @override
  String get progress => 'Progress';

  @override
  String get gems => 'Gems';

  @override
  String get level => 'Level';

  @override
  String get correct => 'Correct!';

  @override
  String get keepPracticing => 'Keep practicing';

  @override
  String get loading => 'Loading...';

  @override
  String get signOut => 'Sign out';

  @override
  String get answerPlaceholder => 'Write your answer...';

  @override
  String get days => 'days';

  @override
  String get activeToday => 'Active today';

  @override
  String get notActiveToday => 'Not active yet today';
}
