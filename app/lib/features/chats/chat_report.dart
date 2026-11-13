import 'reported_message.dart';

class ChatReport {
  const ChatReport({
    required this.id,
    required this.chatId,
    required this.chatTitle,
    required this.reason,
    required this.messages,
    required this.createdAt,
    required this.resolved,
  });

  factory ChatReport.fromJson(Map<String, Object?> json) => ChatReport(
    id: json['id']! as String,
    chatId: json['chatId'] as String?,
    chatTitle: json['chatTitle'] as String? ?? '',
    reason: json['reason'] as String? ?? '',
    messages: [
      for (final item in json['messages'] as List<Object?>? ?? const [])
        ReportedMessage.fromJson(item! as Map<String, Object?>),
    ],
    createdAt: DateTime.parse(json['createdAt']! as String),
    resolved: json['resolvedAt'] != null,
  );

  final String id;
  final String? chatId;
  final String chatTitle;
  final String reason;
  final List<ReportedMessage> messages;
  final DateTime createdAt;
  final bool resolved;
}
