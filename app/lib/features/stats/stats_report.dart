import 'unit_stats.dart';

class StatsReport {
  const StatsReport({
    required this.days,
    required this.threshold,
    required this.full,
    required this.units,
  });

  factory StatsReport.fromJson(Map<String, Object?> json) => StatsReport(
    days: json['days'] as int? ?? 30,
    threshold: json['threshold'] as int? ?? 5,
    full: json['full'] as bool? ?? false,
    units: [
      for (final item in json['units'] as List<Object?>? ?? const [])
        UnitStats.fromJson(item! as Map<String, Object?>),
    ],
  );

  final int days;
  final int threshold;
  final bool full;
  final List<UnitStats> units;
}
