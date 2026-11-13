class Division {
  const Division({
    required this.name,
    required this.from,
    required this.gpuHours,
  });

  factory Division.fromJson(Map<String, Object?> json) => Division(
    name: json['name'] as String? ?? '',
    from: (json['from'] as num?)?.toInt() ?? 0,
    gpuHours: (json['gpuHours'] as num?)?.toInt() ?? 0,
  );

  final String name;
  final int from;
  final int gpuHours;

  Map<String, Object?> toJson() => {
    'name': name,
    'from': from,
    'gpuHours': gpuHours,
  };
}
