import 'run_kind.dart';

class RunInput {
  const RunInput({
    required this.kind,
    required this.title,
    required this.templateId,
    required this.classId,
    this.command = '',
    this.maxHours,
    this.workKind,
    this.workId,
    this.teamId,
  });

  final RunKind kind;
  final String title;
  final String templateId;
  final String classId;
  final String command;
  final double? maxHours;
  final String? workKind;
  final String? workId;
  final String? teamId;

  Map<String, Object?> toJson() => {
    'kind': kind.name,
    'title': title,
    'templateId': templateId,
    'classId': classId,
    if (command.isNotEmpty) 'command': command,
    if (maxHours != null) 'maxHours': maxHours,
    if (workKind != null) 'workKind': workKind,
    if (workId != null) 'workId': workId,
    if (teamId != null) 'teamId': teamId,
  };
}
