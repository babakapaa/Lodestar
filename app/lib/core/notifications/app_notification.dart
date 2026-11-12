class AppNotification {
  const AppNotification({
    required this.id,
    required this.type,
    required this.params,
    required this.link,
    required this.createdAt,
    required this.readAt,
  });

  factory AppNotification.fromJson(Map<String, Object?> json) {
    final params = json['params'];
    final readAt = json['readAt'];
    return AppNotification(
      id: json['id']! as String,
      type: json['type']! as String,
      params: params is Map<String, Object?> ? params : const {},
      link: json['link'] as String? ?? '',
      createdAt: DateTime.parse(json['createdAt']! as String),
      readAt: readAt is String ? DateTime.parse(readAt) : null,
    );
  }

  final String id;
  final String type;
  final Map<String, Object?> params;
  final String link;
  final DateTime createdAt;
  final DateTime? readAt;

  bool get read => readAt != null;

  String? get route =>
      link.startsWith('/') && !link.startsWith('//') ? link : null;

  AppNotification markedRead(DateTime at) => AppNotification(
    id: id,
    type: type,
    params: params,
    link: link,
    createdAt: createdAt,
    readAt: readAt ?? at,
  );
}
