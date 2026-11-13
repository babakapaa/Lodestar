import '../../l10n/generated/app_localizations.dart';
import 'defense_status.dart';
import 'question_kind.dart';

class DefenseLabels {
  const DefenseLabels(this._l10n);

  final AppLocalizations _l10n;

  String kind(QuestionKind kind) => switch (kind) {
    QuestionKind.understanding => _l10n.questionUnderstanding,
    QuestionKind.analysis => _l10n.questionAnalysis,
    QuestionKind.application => _l10n.questionApplication,
    QuestionKind.path => _l10n.questionPath,
  };

  String example(QuestionKind kind) => switch (kind) {
    QuestionKind.understanding => _l10n.questionUnderstandingExample,
    QuestionKind.analysis => _l10n.questionAnalysisExample,
    QuestionKind.application => _l10n.questionApplicationExample,
    QuestionKind.path => _l10n.questionPathExample,
  };

  String mark(int? mark) => switch (mark) {
    0 => _l10n.markNone,
    1 => _l10n.markPartial,
    2 => _l10n.markUnderstands,
    3 => _l10n.markDeep,
    _ => _l10n.markUnset,
  };

  String status(DefenseStatus status) => switch (status) {
    DefenseStatus.planned => _l10n.defensePlanned,
    DefenseStatus.held => _l10n.defenseHeld,
    DefenseStatus.missed => _l10n.defenseMissed,
  };
}
