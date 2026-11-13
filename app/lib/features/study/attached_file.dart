class AttachedFile {
  const AttachedFile({
    required this.id,
    required this.name,
    required this.size,
  });

  factory AttachedFile.fromJson(Map<String, Object?> json) {
    return AttachedFile(
      id: json['id']! as String,
      name: json['name']! as String,
      size: (json['size']! as num).toInt(),
    );
  }

  static List<AttachedFile> list(Object? json) => [
    for (final item in json as List<Object?>? ?? const [])
      AttachedFile.fromJson(item! as Map<String, Object?>),
  ];

  final String id;
  final String name;
  final int size;

  Map<String, Object?> toJson() => {'id': id, 'name': name, 'size': size};
}
