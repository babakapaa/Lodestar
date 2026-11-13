class CourseMaterial {
  const CourseMaterial({
    required this.id,
    required this.title,
    required this.reference,
    required this.passages,
    required this.indexed,
  });

  factory CourseMaterial.fromJson(Map<String, Object?> json) => CourseMaterial(
    id: json['id']! as String,
    title: json['title']! as String,
    reference: json['reference'] as String? ?? '',
    passages: json['passages']! as int,
    indexed: json['indexed']! as int,
  );

  final String id;
  final String title;
  final String reference;
  final int passages;
  final int indexed;

  bool get ready => passages > 0 && indexed >= passages;
}
