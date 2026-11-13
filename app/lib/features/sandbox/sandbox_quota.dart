import 'quota_line.dart';
import 'sandbox_json.dart';
import 'sandbox_term.dart';
import 'team_pool.dart';

class SandboxQuota {
  const SandboxQuota({
    required this.term,
    required this.lines,
    required this.hours,
    required this.reserved,
    required this.spent,
    required this.forecast,
    required this.teams,
  });

  factory SandboxQuota.fromJson(Map<String, Object?> json) => SandboxQuota(
    term: json['term'] == null
        ? null
        : SandboxTerm.fromJson(json['term']! as Map<String, Object?>),
    lines: SandboxJson.list(json['lines'], QuotaLine.fromJson),
    hours: SandboxJson.number(json['hours']),
    reserved: SandboxJson.number(json['reserved']),
    spent: SandboxJson.number(json['spent']),
    forecast: SandboxJson.number(json['forecast']),
    teams: SandboxJson.list(json['teams'], TeamPool.fromJson),
  );

  final SandboxTerm? term;
  final List<QuotaLine> lines;
  final double hours;
  final double reserved;
  final double spent;
  final double forecast;
  final List<TeamPool> teams;

  double get granted =>
      lines.fold(0, (sum, l) => l.hours > 0 ? sum + l.hours : sum);

  double get available => hours - reserved;

  double get share => granted <= 0 ? 0 : (hours / granted).clamp(0, 1);
}
