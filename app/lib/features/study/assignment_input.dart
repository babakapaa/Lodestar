import 'grading.dart';

class AssignmentInput {
  const AssignmentInput({
    required this.title,
    required this.description,
    required this.dueAt,
    required this.grading,
    required this.maxPoints,
    required this.published,
    required this.files,
  });

  final String title;
  final String description;
  final DateTime? dueAt;
  final Grading grading;
  final double? maxPoints;
  final bool published;
  final List<String> files;

  Map<String, Object?> toJson() => {
    'title': title,
    'description': description,
    'dueAt': dueAt?.toUtc().toIso8601String(),
    'grading': grading.name,
    'maxPoints': grading == Grading.points ? maxPoints : null,
    'published': published,
    'files': files,
  };
}
