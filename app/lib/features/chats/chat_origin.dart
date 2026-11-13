class ChatOrigin {
  const ChatOrigin({required this.kind, required this.id});

  factory ChatOrigin.fromJson(Map<String, Object?> json) =>
      ChatOrigin(kind: json['kind']! as String, id: json['id']! as String);

  static const course = 'study.course';

  final String kind;
  final String id;

  Map<String, Object?> toJson() => {'kind': kind, 'id': id};
}
