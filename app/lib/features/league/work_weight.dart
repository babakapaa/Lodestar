class WorkWeight {
  const WorkWeight({
    required this.kind,
    required this.id,
    required this.title,
    required this.project,
    required this.weight,
  });

  factory WorkWeight.fromJson(Map<String, Object?> json) => WorkWeight(
    kind: json['kind']! as String,
    id: json['id']! as String,
    title: json['title'] as String? ?? '',
    project: json['project'] as bool? ?? false,
    weight: (json['weight'] as num?)?.toInt() ?? 0,
  );

  final String kind;
  final String id;
  final String title;
  final bool project;
  final int weight;

  Map<String, Object?> toJson([int? value]) => {
    'kind': kind,
    'id': id,
    'weight': value ?? weight,
  };
}
