class ReportedMessage {
  const ReportedMessage({
    required this.seq,
    required this.authorName,
    required this.text,
    required this.createdAt,
  });

  factory ReportedMessage.fromJson(Map<String, Object?> json) =>
      ReportedMessage(
        seq: (json['seq']! as num).toInt(),
        authorName: json['authorName'] as String? ?? '',
        text: json['text'] as String? ?? '',
        createdAt: DateTime.parse(json['createdAt']! as String),
      );

  final int seq;
  final String authorName;
  final String text;
  final DateTime createdAt;
}
