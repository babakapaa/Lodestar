import 'chat_message.dart';

class MessagePage {
  const MessagePage({required this.messages, required this.more});

  factory MessagePage.fromJson(Map<String, Object?> json) => MessagePage(
    messages: [
      for (final item in json['messages']! as List<Object?>)
        ChatMessage.fromJson(item! as Map<String, Object?>),
    ],
    more: json['more']! as bool,
  );

  final List<ChatMessage> messages;
  final bool more;
}
