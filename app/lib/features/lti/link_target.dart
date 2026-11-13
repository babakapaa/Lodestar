class LinkTarget {
  const LinkTarget({
    required this.kind,
    required this.id,
    required this.title,
    required this.path,
  });

  factory LinkTarget.fromJson(Map<String, Object?> json) => LinkTarget(
    kind: json['kind']! as String,
    id: json['id']! as String,
    title: json['title'] as String? ?? '',
    path: json['path'] as String? ?? '',
  );

  final String kind;
  final String id;
  final String title;
  final String path;
}
