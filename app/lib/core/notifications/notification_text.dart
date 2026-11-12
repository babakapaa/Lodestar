class NotificationText {
  const NotificationText({required this.title, required this.body});

  factory NotificationText.fromJson(Map<String, Object?> json) {
    return NotificationText(
      title: json['title']! as String,
      body: json['body']! as String,
    );
  }

  final String title;
  final String body;
}
