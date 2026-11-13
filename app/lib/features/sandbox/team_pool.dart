import 'quota_line.dart';
import 'sandbox_json.dart';

class TeamPool {
  const TeamPool({
    required this.id,
    required this.name,
    required this.hours,
    required this.reserved,
    required this.lines,
  });

  factory TeamPool.fromJson(Map<String, Object?> json) => TeamPool(
    id: json['id']! as String,
    name: json['name'] as String? ?? '',
    hours: SandboxJson.number(json['hours']),
    reserved: SandboxJson.number(json['reserved']),
    lines: SandboxJson.list(json['lines'], QuotaLine.fromJson),
  );

  final String id;
  final String name;
  final double hours;
  final double reserved;
  final List<QuotaLine> lines;

  double get available => hours - reserved;
}
