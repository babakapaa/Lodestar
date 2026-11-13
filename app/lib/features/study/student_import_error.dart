class StudentImportError {
  const StudentImportError({required this.line, required this.code});

  factory StudentImportError.fromJson(Map<String, Object?> json) =>
      StudentImportError(
        line: (json['line'] as num?)?.toInt() ?? 0,
        code: json['code'] as String? ?? '',
      );

  final int line;
  final String code;
}
