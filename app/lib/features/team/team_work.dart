class TeamWork {
  const TeamWork({required this.text, required this.links});

  factory TeamWork.fromJson(Object? json) {
    final map = json is Map<String, Object?> ? json : const <String, Object?>{};
    return TeamWork(
      text: (map['body'] ?? map['text']) as String? ?? '',
      links: [
        for (final link in map['links'] as List<Object?>? ?? const [])
          link! as String,
      ],
    );
  }

  final String text;
  final List<String> links;
}
