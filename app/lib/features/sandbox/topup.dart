import 'sandbox_json.dart';

class Topup {
  const Topup({
    required this.id,
    required this.userId,
    required this.displayName,
    required this.workKind,
    required this.workId,
    required this.title,
    required this.hours,
    required this.reason,
    required this.status,
    required this.granted,
    required this.createdAt,
    required this.mine,
  });

  factory Topup.fromJson(Map<String, Object?> json) => Topup(
    id: json['id']! as String,
    userId: json['userId'] as String? ?? '',
    displayName:
        json['displayName'] as String? ?? json['email'] as String? ?? '',
    workKind: json['workKind'] as String? ?? '',
    workId: json['workId'] as String? ?? '',
    title: json['title'] as String? ?? '',
    hours: SandboxJson.number(json['hours']),
    reason: json['reason'] as String? ?? '',
    status: json['status'] as String? ?? 'pending',
    granted: SandboxJson.maybeNumber(json['granted']),
    createdAt: SandboxJson.time(json['createdAt']) ?? DateTime.now(),
    mine: json['mine'] as bool? ?? false,
  );

  final String id;
  final String userId;
  final String displayName;
  final String workKind;
  final String workId;
  final String title;
  final double hours;
  final String reason;
  final String status;
  final double? granted;
  final DateTime createdAt;
  final bool mine;

  bool get pending => status == 'pending';
}
