class MapRules {
  const MapRules({required this.minSources, required this.minScholarly});

  factory MapRules.fromJson(Map<String, Object?> json) => MapRules(
    minSources: (json['minSources'] as num?)?.toInt() ?? 0,
    minScholarly: (json['minScholarly'] as num?)?.toInt() ?? 0,
  );

  static MapRules? maybe(Object? json) =>
      json is Map<String, Object?> ? MapRules.fromJson(json) : null;

  static const none = MapRules(minSources: 0, minScholarly: 0);

  final int minSources;
  final int minScholarly;

  Map<String, Object?> toJson() => {
    'minSources': minSources,
    'minScholarly': minScholarly,
  };
}
