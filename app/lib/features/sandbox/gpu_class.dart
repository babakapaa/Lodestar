import 'sandbox_json.dart';

class GpuClass {
  const GpuClass({
    required this.id,
    required this.name,
    required this.memoryGb,
    required this.rate,
    required this.archived,
  });

  factory GpuClass.fromJson(Map<String, Object?> json) => GpuClass(
    id: json['id']! as String,
    name: json['name'] as String? ?? '',
    memoryGb: json['memoryGb'] as int? ?? 0,
    rate: SandboxJson.number(json['rate']),
    archived: json['archived'] as bool? ?? false,
  );

  final String id;
  final String name;
  final int memoryGb;
  final double rate;
  final bool archived;

  Map<String, Object?> toJson() => {
    'name': name,
    'memoryGb': memoryGb,
    'rate': rate,
    'archived': archived,
  };
}
