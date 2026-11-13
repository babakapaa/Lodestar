import 'chat_file.dart';

class ChatMessage {
  const ChatMessage({
    required this.id,
    required this.chatId,
    required this.seq,
    required this.rev,
    required this.authorId,
    required this.clientId,
    required this.text,
    required this.replyTo,
    required this.mentions,
    required this.files,
    required this.views,
    required this.deleted,
    required this.editedAt,
    required this.createdAt,
  });

  factory ChatMessage.fromJson(Map<String, Object?> json) => ChatMessage(
    id: json['id']! as String,
    chatId: json['chatId']! as String,
    seq: (json['seq']! as num).toInt(),
    rev: (json['rev']! as num).toInt(),
    authorId: json['authorId'] as String?,
    clientId: json['clientId'] as String?,
    text: json['text'] as String? ?? '',
    replyTo: (json['replyTo'] as num?)?.toInt(),
    mentions: [
      for (final id in json['mentions'] as List<Object?>? ?? const [])
        id! as String,
    ],
    files: [
      for (final item in json['files'] as List<Object?>? ?? const [])
        ChatFile.fromJson(item! as Map<String, Object?>),
    ],
    views: (json['views'] as num?)?.toInt(),
    deleted: json['deleted'] as bool? ?? false,
    editedAt: _time(json['editedAt']),
    createdAt: DateTime.parse(json['createdAt']! as String),
  );

  final String id;
  final String chatId;
  final int seq;
  final int rev;
  final String? authorId;
  final String? clientId;
  final String text;
  final int? replyTo;
  final List<String> mentions;
  final List<ChatFile> files;
  final int? views;
  final bool deleted;
  final DateTime? editedAt;
  final DateTime createdAt;

  Map<String, Object?> toJson() => {
    'id': id,
    'chatId': chatId,
    'seq': seq,
    'rev': rev,
    'authorId': authorId,
    'clientId': clientId,
    'text': text,
    'replyTo': replyTo,
    'mentions': mentions,
    'files': [for (final file in files) file.toJson()],
    'views': views,
    'deleted': deleted,
    'editedAt': editedAt?.toIso8601String(),
    'createdAt': createdAt.toIso8601String(),
  };

  static DateTime? _time(Object? value) =>
      value == null ? null : DateTime.parse(value as String);
}
