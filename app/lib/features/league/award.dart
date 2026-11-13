class Award {
  const Award({
    required this.id,
    required this.title,
    required this.about,
    required this.gpuHours,
    required this.quantity,
    required this.granted,
    required this.archived,
  });

  factory Award.fromJson(Map<String, Object?> json) => Award(
    id: json['id']! as String,
    title: json['title'] as String? ?? '',
    about: json['about'] as String? ?? '',
    gpuHours: (json['gpuHours'] as num?)?.toInt() ?? 0,
    quantity: (json['quantity'] as num?)?.toInt(),
    granted: (json['granted'] as num?)?.toInt() ?? 0,
    archived: json['archivedAt'] != null,
  );

  final String id;
  final String title;
  final String about;
  final int gpuHours;
  final int? quantity;
  final int granted;
  final bool archived;

  bool get left => !archived && (quantity == null || granted < quantity!);
}
