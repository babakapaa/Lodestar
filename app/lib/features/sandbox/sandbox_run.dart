import 'run_kind.dart';
import 'run_status.dart';
import 'sandbox_json.dart';

class SandboxRun {
  const SandboxRun({
    required this.id,
    required this.userId,
    required this.displayName,
    required this.teamId,
    required this.teamName,
    required this.kind,
    required this.title,
    required this.command,
    required this.templateId,
    required this.templateName,
    required this.classId,
    required this.className,
    required this.rate,
    required this.maxHours,
    required this.workKind,
    required this.workId,
    required this.workTitle,
    required this.status,
    required this.position,
    required this.url,
    required this.exitCode,
    required this.reason,
    required this.charged,
    required this.queuedAt,
    required this.startedAt,
    required this.endedAt,
    required this.heldAt,
    required this.holdNote,
    required this.refunded,
    required this.gpuLoad,
    required this.log,
  });

  factory SandboxRun.fromJson(Map<String, Object?> json) => SandboxRun(
    id: json['id']! as String,
    userId: json['userId'] as String? ?? '',
    displayName: json['displayName'] as String? ?? '',
    teamId: json['teamId'] as String?,
    teamName: json['teamName'] as String? ?? '',
    kind: RunKind.parse(json['kind'] as String?),
    title: json['title'] as String? ?? '',
    command: json['command'] as String? ?? '',
    templateId: json['templateId'] as String? ?? '',
    templateName: json['templateName'] as String? ?? '',
    classId: json['classId'] as String? ?? '',
    className: json['className'] as String? ?? '',
    rate: SandboxJson.number(json['rate']),
    maxHours: SandboxJson.number(json['maxHours']),
    workKind: json['workKind'] as String?,
    workId: json['workId'] as String?,
    workTitle: json['workTitle'] as String? ?? '',
    status: RunStatus.parse(json['status'] as String?),
    position: json['position'] as int?,
    url: json['url'] as String? ?? '',
    exitCode: json['exitCode'] as int?,
    reason: json['reason'] as String? ?? '',
    charged: SandboxJson.maybeNumber(json['charged']),
    queuedAt: SandboxJson.time(json['queuedAt']) ?? DateTime.now(),
    startedAt: SandboxJson.time(json['startedAt']),
    endedAt: SandboxJson.time(json['endedAt']),
    heldAt: SandboxJson.time(json['heldAt']),
    holdNote: json['holdNote'] as String? ?? '',
    refunded: json['refunded'] as bool? ?? false,
    gpuLoad: SandboxJson.maybeNumber(json['gpuLoad']),
    log: json['log'] as String? ?? '',
  );

  final String id;
  final String userId;
  final String displayName;
  final String? teamId;
  final String teamName;
  final RunKind kind;
  final String title;
  final String command;
  final String templateId;
  final String templateName;
  final String classId;
  final String className;
  final double rate;
  final double maxHours;
  final String? workKind;
  final String? workId;
  final String workTitle;
  final RunStatus status;
  final int? position;
  final String url;
  final int? exitCode;
  final String reason;
  final double? charged;
  final DateTime queuedAt;
  final DateTime? startedAt;
  final DateTime? endedAt;
  final DateTime? heldAt;
  final String holdNote;
  final bool refunded;
  final double? gpuLoad;
  final String log;

  bool get free => rate == 0;

  bool get wasteful =>
      !status.active && gpuLoad != null && gpuLoad! < 30 && !free;
}
