class TaskInput {
  const TaskInput({
    required this.title,
    required this.about,
    required this.partner,
    required this.points,
    required this.minSize,
    required this.maxSize,
    required this.selfFormed,
    required this.dueAt,
    this.closed = false,
  });

  final String title;
  final String about;
  final String partner;
  final int points;
  final int minSize;
  final int maxSize;
  final bool selfFormed;
  final DateTime? dueAt;
  final bool closed;

  Map<String, Object?> toJson() => {
    'title': title,
    'about': about,
    'partner': partner,
    'points': points,
    'minSize': minSize,
    'maxSize': maxSize,
    'selfFormed': selfFormed,
    'dueAt': dueAt?.toUtc().toIso8601String(),
    'closed': closed,
  };
}
