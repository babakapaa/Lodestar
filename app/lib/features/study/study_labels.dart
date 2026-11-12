import 'package:intl/intl.dart';

import '../../l10n/generated/app_localizations.dart';
import 'assignment.dart';
import 'course_role.dart';
import 'grade.dart';
import 'grading.dart';

class StudyLabels {
  const StudyLabels(this._l10n, this._locale);

  final AppLocalizations _l10n;
  final String _locale;

  String role(CourseRole role) => switch (role) {
    CourseRole.teacher => _l10n.studyRoleTeacher,
    CourseRole.student => _l10n.studyRoleStudent,
    CourseRole.manager => _l10n.studyRoleManager,
  };

  String grading(Grading grading) => switch (grading) {
    Grading.points => _l10n.gradingPoints,
    Grading.pass => _l10n.gradingPass,
    Grading.none => _l10n.gradingNone,
  };

  String time(DateTime at) =>
      DateFormat.MMMd(_locale).add_Hm().format(at.toLocal());

  String due(DateTime? dueAt) =>
      dueAt == null ? _l10n.assignmentNoDue : _l10n.assignmentDue(time(dueAt));

  String grade(Grade grade, Assignment assignment) {
    final points = grade.points;
    final passed = grade.passed;
    if (points != null) {
      return _l10n.gradePoints(
        _number(points),
        _number(assignment.maxPoints ?? 0),
      );
    }
    if (passed != null) {
      return passed ? _l10n.gradePassed : _l10n.gradeFailed;
    }
    return _l10n.gradeReviewed;
  }

  String _number(double value) =>
      NumberFormat.decimalPattern(_locale).format(value);
}
