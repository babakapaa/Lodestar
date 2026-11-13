class NewSlots {
  const NewSlots({
    required this.startsAt,
    required this.minutes,
    required this.place,
    required this.capacity,
    required this.count,
  });

  final DateTime startsAt;
  final int minutes;
  final String place;
  final int capacity;
  final int count;

  Map<String, Object?> toJson() => {
    'startsAt': startsAt.toUtc().toIso8601String(),
    'minutes': minutes,
    'place': place,
    'capacity': capacity,
    'count': count,
  };
}
