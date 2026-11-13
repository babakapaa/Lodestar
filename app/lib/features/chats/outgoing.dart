class Outgoing {
  const Outgoing({
    required this.clientId,
    required this.chatId,
    required this.text,
    required this.replyTo,
    required this.createdAt,
    required this.error,
  });

  final String clientId;
  final String chatId;
  final String text;
  final int? replyTo;
  final DateTime createdAt;
  final String? error;

  Map<String, Object?> toJson() => {
    'clientId': clientId,
    'text': text,
    'replyTo': ?replyTo,
  };
}
