import 'attendance_counts.dart';

class JournalStudent {
  const JournalStudent({
    required this.userId,
    required this.name,
    required this.counts,
    required this.rate,
    required this.streak,
    required this.longAbsent,
  });

  factory JournalStudent.fromJson(Map<String, Object?> json) => JournalStudent(
    userId: json['userId']! as String,
    name: json['name'] as String? ?? '',
    counts: AttendanceCounts.fromJson(json['counts'] as Map<String, Object?>?),
    rate: (json['rate'] as num?)?.toDouble(),
    streak: json['streak'] as int? ?? 0,
    longAbsent: json['longAbsent'] as bool? ?? false,
  );

  final String userId;
  final String name;
  final AttendanceCounts counts;
  final double? rate;
  final int streak;
  final bool longAbsent;
}
