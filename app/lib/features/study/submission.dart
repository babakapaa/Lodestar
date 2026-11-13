import 'attached_file.dart';
import 'search_map.dart';

class Submission {
  const Submission({
    required this.id,
    required this.body,
    required this.submittedAt,
    required this.files,
    this.searchMap,
  });

  factory Submission.fromJson(Map<String, Object?> json) {
    return Submission(
      id: json['id']! as String,
      body: json['body'] as String? ?? '',
      submittedAt: DateTime.parse(json['submittedAt']! as String),
      files: AttachedFile.list(json['files']),
      searchMap: SearchMap.maybe(json['searchMap']),
    );
  }

  static Submission? maybe(Object? json) =>
      json is Map<String, Object?> ? Submission.fromJson(json) : null;

  final String id;
  final String body;
  final DateTime submittedAt;
  final List<AttachedFile> files;
  final SearchMap? searchMap;
}
