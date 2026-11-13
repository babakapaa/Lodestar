import 'chat_file.dart';

class Outgoing {
  const Outgoing({
    required this.clientId,
    required this.chatId,
    required this.text,
    required this.replyTo,
    required this.mentions,
    required this.files,
    required this.createdAt,
    required this.error,
  });

  final String clientId;
  final String chatId;
  final String text;
  final int? replyTo;
  final List<String> mentions;
  final List<ChatFile> files;
  final DateTime createdAt;
  final String? error;

  Map<String, Object?> toJson() => {
    'clientId': clientId,
    'text': text,
    'replyTo': ?replyTo,
    if (mentions.isNotEmpty) 'mentions': mentions,
    if (files.isNotEmpty) 'fileIds': [for (final file in files) file.id],
  };
}
