class StatsEntry {
  const StatsEntry({required this.unitId, required this.name});

  factory StatsEntry.fromJson(Map<String, Object?> json) => StatsEntry(
    unitId: json['unitId']! as String,
    name: json['name']! as String,
  );

  final String unitId;
  final String name;
}
