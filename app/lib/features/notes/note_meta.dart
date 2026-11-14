import 'package:flutter/foundation.dart';

@immutable
class NoteMeta {
  const NoteMeta({
    required this.path,
    this.tags = const [],
    this.links = const [],
    this.assignmentId,
    this.lessonId,
  });

  factory NoteMeta.fromJson(Map<String, Object?> json) => NoteMeta(
    path: json['path']! as String,
    tags: _strings(json['tags']),
    links: _strings(json['links']),
    assignmentId: json['assignmentId'] as String?,
    lessonId: json['lessonId'] as String?,
  );

  final String path;
  final List<String> tags;
  final List<String> links;
  final String? assignmentId;
  final String? lessonId;

  static const canvasExtension = '.canvas';

  String get name => path.split('/').last;

  bool get isCanvas => path.toLowerCase().endsWith(canvasExtension);

  String get label =>
      isCanvas ? name.substring(0, name.length - canvasExtension.length) : name;

  String get folder {
    final at = path.lastIndexOf('/');
    return at < 0 ? '' : path.substring(0, at);
  }

  NoteMeta copyWith({String? path, List<String>? tags, List<String>? links}) =>
      NoteMeta(
        path: path ?? this.path,
        tags: tags ?? this.tags,
        links: links ?? this.links,
        assignmentId: assignmentId,
        lessonId: lessonId,
      );

  Map<String, Object?> toJson() => {
    'path': path,
    'tags': tags,
    'links': links,
    'assignmentId': assignmentId,
    'lessonId': lessonId,
  };

  static String normalizePath(String path) => [
    for (final part in path.split('/'))
      if (part.trim().isNotEmpty) part.trim(),
  ].join('/');

  static List<String> _strings(Object? value) => [
    for (final item in (value as List<Object?>?) ?? const []) item! as String,
  ];
}
