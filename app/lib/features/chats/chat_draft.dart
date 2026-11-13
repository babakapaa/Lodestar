import 'chat_file.dart';

class ChatDraft {
  const ChatDraft({
    required this.text,
    required this.mentions,
    required this.files,
  });

  final String text;
  final List<String> mentions;
  final List<ChatFile> files;
}
