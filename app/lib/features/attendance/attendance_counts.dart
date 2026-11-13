class AttendanceCounts {
  const AttendanceCounts({
    required this.present,
    required this.late,
    required this.absent,
    required this.excused,
    required this.unmarked,
  });

  factory AttendanceCounts.fromJson(Map<String, Object?>? json) =>
      AttendanceCounts(
        present: json?['present'] as int? ?? 0,
        late: json?['late'] as int? ?? 0,
        absent: json?['absent'] as int? ?? 0,
        excused: json?['excused'] as int? ?? 0,
        unmarked: json?['unmarked'] as int? ?? 0,
      );

  final int present;
  final int late;
  final int absent;
  final int excused;
  final int unmarked;

  int get here => present + late;
}
