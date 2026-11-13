class Call {
  const Call({
    required this.id,
    required this.chatId,
    required this.startedBy,
    required this.startedAt,
    required this.people,
  });

  factory Call.fromJson(Map<String, Object?> json) => Call(
    id: json['id']! as String,
    chatId: json['chatId']! as String,
    startedBy: json['startedBy'] as String?,
    startedAt: DateTime.parse(json['startedAt']! as String),
    people: [
      for (final item in json['people'] as List<Object?>? ?? const [])
        item! as String,
    ],
  );

  final String id;
  final String chatId;
  final String? startedBy;
  final DateTime startedAt;
  final List<String> people;
}
