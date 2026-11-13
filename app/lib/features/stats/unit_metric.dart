class UnitMetric {
  const UnitMetric({required this.key, required this.value});

  factory UnitMetric.fromJson(Map<String, Object?> json) => UnitMetric(
    key: json['key']! as String,
    value: (json['value'] as num? ?? 0).toDouble(),
  );

  final String key;
  final double value;
}
