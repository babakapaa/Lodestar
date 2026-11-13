import 'student_import_error.dart';

class StudentImport {
  const StudentImport({
    required this.added,
    required this.existing,
    required this.errors,
  });

  factory StudentImport.fromJson(Map<String, Object?> json) => StudentImport(
    added: (json['added'] as num?)?.toInt() ?? 0,
    existing: (json['existing'] as num?)?.toInt() ?? 0,
    errors: [
      for (final item in json['errors'] as List<Object?>? ?? const [])
        StudentImportError.fromJson(item! as Map<String, Object?>),
    ],
  );

  final int added;
  final int existing;
  final List<StudentImportError> errors;
}
