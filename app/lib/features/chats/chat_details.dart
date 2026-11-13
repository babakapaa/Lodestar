import 'chat.dart';
import 'chat_member.dart';

class ChatDetails {
  const ChatDetails({
    required this.chat,
    required this.members,
    required this.peerReadSeq,
  });

  factory ChatDetails.fromJson(Map<String, Object?> json) => ChatDetails(
    chat: Chat.fromJson(json),
    members: [
      for (final item in json['members'] as List<Object?>? ?? const [])
        ChatMember.fromJson(item! as Map<String, Object?>),
    ],
    peerReadSeq: (json['peerReadSeq'] as num?)?.toInt(),
  );

  final Chat chat;
  final List<ChatMember> members;
  final int? peerReadSeq;
}
