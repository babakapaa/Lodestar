import 'chat_role.dart';

class ChatMember {
  const ChatMember({
    required this.userId,
    required this.displayName,
    required this.role,
    required this.synced,
  });

  factory ChatMember.fromJson(Map<String, Object?> json) => ChatMember(
    userId: json['userId']! as String,
    displayName: json['displayName'] as String? ?? '',
    role: ChatRole.parse(json['role']! as String),
    synced: json['synced'] as bool? ?? false,
  );

  final String userId;
  final String displayName;
  final ChatRole role;
  final bool synced;
}
