import 'unit_metric.dart';

class UnitStats {
  const UnitStats({
    required this.unitId,
    required this.kind,
    required this.name,
    required this.people,
    required this.hidden,
    required this.metrics,
  });

  factory UnitStats.fromJson(Map<String, Object?> json) => UnitStats(
    unitId: json['unitId']! as String,
    kind: json['kind'] as String? ?? '',
    name: json['name'] as String? ?? '',
    people: json['people'] as int?,
    hidden: json['hidden'] as bool? ?? false,
    metrics: [
      for (final item in json['metrics'] as List<Object?>? ?? const [])
        UnitMetric.fromJson(item! as Map<String, Object?>),
    ],
  );

  final String unitId;
  final String kind;
  final String name;
  final int? people;
  final bool hidden;
  final List<UnitMetric> metrics;
}
