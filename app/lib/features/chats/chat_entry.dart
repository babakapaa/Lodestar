import 'chat_message.dart';
import 'outgoing.dart';

class ChatEntry {
  const ChatEntry.sent(ChatMessage this.message) : outgoing = null;

  const ChatEntry.pending(Outgoing this.outgoing) : message = null;

  final ChatMessage? message;
  final Outgoing? outgoing;

  String get key => message?.id ?? outgoing!.clientId;
}
