class AwardInput {
  const AwardInput({
    required this.title,
    required this.about,
    required this.gpuHours,
    required this.quantity,
    this.archived = false,
  });

  final String title;
  final String about;
  final int gpuHours;
  final int? quantity;
  final bool archived;

  Map<String, Object?> toJson() => {
    'title': title,
    'about': about,
    'gpuHours': gpuHours,
    'quantity': quantity,
    'archived': archived,
  };
}
