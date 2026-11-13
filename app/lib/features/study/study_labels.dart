import 'package:intl/intl.dart';

import '../../l10n/generated/app_localizations.dart';
import 'assignment.dart';
import 'course_role.dart';
import 'defense_format.dart';
import 'grade.dart';
import 'grading.dart';
import 'lesson_status.dart';
import 'navigator_mode.dart';

class StudyLabels {
  const StudyLabels(this._l10n, this._locale);

  final AppLocalizations _l10n;
  final String _locale;

  String role(CourseRole role) => switch (role) {
    CourseRole.teacher => _l10n.studyRoleTeacher,
    CourseRole.student => _l10n.studyRoleStudent,
    CourseRole.manager => _l10n.studyRoleManager,
  };

  String navigatorMode(NavigatorMode mode) => switch (mode) {
    NavigatorMode.course => _l10n.navigatorModeCourse,
    NavigatorMode.study => _l10n.navigatorModeStudyShort,
    NavigatorMode.project => _l10n.navigatorModeProjectShort,
    NavigatorMode.closed => _l10n.navigatorModeClosedShort,
  };

  String defenseFormat(DefenseFormat format) => switch (format) {
    DefenseFormat.none => _l10n.defenseFormatNone,
    DefenseFormat.full => _l10n.defenseFormatFull,
    DefenseFormat.short => _l10n.defenseFormatShort,
    DefenseFormat.sample => _l10n.defenseFormatSample,
  };

  String grading(Grading grading) => switch (grading) {
    Grading.points => _l10n.gradingPoints,
    Grading.pass => _l10n.gradingPass,
    Grading.none => _l10n.gradingNone,
  };

  String kind(String kind) => switch (kind) {
    'lecture' => _l10n.lessonKindLecture,
    'practice' => _l10n.lessonKindPractice,
    'lab' => _l10n.lessonKindLab,
    'seminar' => _l10n.lessonKindSeminar,
    'exam' => _l10n.lessonKindExam,
    'consultation' => _l10n.lessonKindConsultation,
    _ => kind,
  };

  String? status(LessonStatus status) => switch (status) {
    LessonStatus.planned => null,
    LessonStatus.changed => _l10n.lessonChanged,
    LessonStatus.cancelled => _l10n.lessonCancelled,
  };

  String clock(DateTime at) => DateFormat.Hm(_locale).format(at.toLocal());

  String day(DateTime day) => DateFormat.MMMMEEEEd(_locale).format(day);

  String shortDay(DateTime day) => DateFormat.MMMd(_locale).format(day);

  String date(DateTime day) => DateFormat.yMMMd(_locale).format(day);

  String weekday(int day) =>
      DateFormat.EEEE(_locale).format(DateTime(2024, 1, day));

  String shortWeekday(int day) =>
      DateFormat.E(_locale).format(DateTime(2024, 1, day));

  String points(double value) => _number(value);

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
