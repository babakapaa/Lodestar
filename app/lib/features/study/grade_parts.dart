class GradeParts {
  const GradeParts({this.defense, this.map, this.work});

  factory GradeParts.fromJson(Map<String, Object?> json) => GradeParts(
    defense: (json['defense'] as num?)?.toDouble(),
    map: (json['map'] as num?)?.toDouble(),
    work: (json['work'] as num?)?.toDouble(),
  );

  static GradeParts? maybe(Object? json) =>
      json is Map<String, Object?> ? GradeParts.fromJson(json) : null;

  final double? defense;
  final double? map;
  final double? work;

  Map<String, Object?> toJson() => {
    'defense': ?defense,
    'map': ?map,
    'work': ?work,
  };
}
