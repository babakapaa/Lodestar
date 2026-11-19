import 'grade_parts.dart';

class GradeWeights {
  const GradeWeights({
    required this.defense,
    required this.map,
    required this.work,
  });

  factory GradeWeights.fromJson(Map<String, Object?> json) => GradeWeights(
    defense: json['defense'] as int? ?? 0,
    map: json['map'] as int? ?? 0,
    work: json['work'] as int? ?? 0,
  );

  static GradeWeights? maybe(Object? json) =>
      json is Map<String, Object?> ? GradeWeights.fromJson(json) : null;

  static const suggested = GradeWeights(defense: 50, map: 30, work: 20);

  final int defense;
  final int map;
  final int work;

  bool get valid =>
      defense >= 0 && map >= 0 && work >= 0 && defense + map + work == 100;

  double? points(GradeParts parts, double maxPoints) {
    var share = 0.0;
    for (final (weight, value) in [
      (defense, parts.defense),
      (map, parts.map),
      (work, parts.work),
    ]) {
      if (weight == 0) {
        continue;
      }
      if (value == null || value < 0 || value > 100) {
        return null;
      }
      share += weight * value;
    }
    return (share / 10000 * maxPoints * 100).roundToDouble() / 100;
  }

  Map<String, Object?> toJson() => {
    'defense': defense,
    'map': map,
    'work': work,
  };
}
