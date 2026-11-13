import 'sandbox_json.dart';

class GpuNode {
  const GpuNode({
    required this.id,
    required this.name,
    required this.classId,
    required this.className,
    required this.gpus,
    required this.busy,
    required this.online,
    required this.seenAt,
    this.token,
  });

  factory GpuNode.fromJson(Map<String, Object?> json) => GpuNode(
    id: json['id']! as String,
    name: json['name'] as String? ?? '',
    classId: json['classId'] as String? ?? '',
    className: json['className'] as String? ?? '',
    gpus: json['gpus'] as int? ?? 0,
    busy: json['busy'] as int? ?? 0,
    online: json['online'] as bool? ?? false,
    seenAt: SandboxJson.time(json['seenAt']),
    token: json['token'] as String?,
  );

  final String id;
  final String name;
  final String classId;
  final String className;
  final int gpus;
  final int busy;
  final bool online;
  final DateTime? seenAt;
  final String? token;
}
