import 'assignment_source.dart';
import 'link_target.dart';

class LtiContext {
  const LtiContext({
    required this.id,
    required this.title,
    required this.label,
    required this.platform,
    required this.instructor,
    required this.link,
    required this.options,
    required this.source,
    required this.sourceEditable,
    required this.canSync,
    required this.syncedAt,
    required this.syncError,
  });

  factory LtiContext.fromJson(Map<String, Object?> json) {
    final link = json['link'] as Map<String, Object?>?;
    final synced = json['syncedAt'] as String?;
    return LtiContext(
      id: json['id']! as String,
      title: json['title'] as String? ?? '',
      label: json['label'] as String? ?? '',
      platform: json['platform'] as String? ?? '',
      instructor: json['instructor'] as bool? ?? false,
      link: link == null ? null : LinkTarget.fromJson(link),
      options: [
        for (final item in json['options'] as List<Object?>? ?? const [])
          LinkTarget.fromJson(item! as Map<String, Object?>),
      ],
      source: AssignmentSource.parse(json['assignmentSource']),
      sourceEditable: json['sourceEditable'] as bool? ?? false,
      canSync: json['canSync'] as bool? ?? false,
      syncedAt: synced == null ? null : DateTime.parse(synced).toLocal(),
      syncError: json['syncError'] as String?,
    );
  }

  final String id;
  final String title;
  final String label;
  final String platform;
  final bool instructor;
  final LinkTarget? link;
  final List<LinkTarget> options;
  final AssignmentSource source;
  final bool sourceEditable;
  final bool canSync;
  final DateTime? syncedAt;
  final String? syncError;

  String get name =>
      label.isEmpty || label == title ? title : '$title ($label)';
}
