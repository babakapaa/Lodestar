class UploadedFile {
  const UploadedFile({
    required this.id,
    required this.name,
    required this.contentType,
    required this.size,
  });

  factory UploadedFile.fromJson(Map<String, Object?> json) => UploadedFile(
    id: json['id']! as String,
    name: json['name']! as String,
    contentType: json['contentType']! as String,
    size: (json['size']! as num).toInt(),
  );

  final String id;
  final String name;
  final String contentType;
  final int size;
}
