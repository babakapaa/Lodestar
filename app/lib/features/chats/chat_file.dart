class ChatFile {
  const ChatFile({
    required this.id,
    required this.name,
    required this.contentType,
    required this.size,
  });

  factory ChatFile.fromJson(Map<String, Object?> json) => ChatFile(
    id: json['id']! as String,
    name: json['name']! as String,
    contentType: json['contentType']! as String,
    size: (json['size']! as num).toInt(),
  );

  final String id;
  final String name;
  final String contentType;
  final int size;

  Map<String, Object?> toJson() => {
    'id': id,
    'name': name,
    'contentType': contentType,
    'size': size,
  };
}
