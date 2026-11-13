class WeekCount {
  const WeekCount({required this.week, required this.students});

  factory WeekCount.fromJson(Map<String, Object?> json) => WeekCount(
    week: DateTime.parse(json['week']! as String),
    students: json['students'] as int?,
  );

  final DateTime week;
  final int? students;
}
