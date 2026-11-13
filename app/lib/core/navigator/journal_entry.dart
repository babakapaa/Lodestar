class JournalEntry {
  const JournalEntry({
    required this.id,
    required this.query,
    required this.sources,
    required this.createdAt,
  });

  factory JournalEntry.fromJson(Map<String, Object?> json) => JournalEntry(
    id: json['id']! as String,
    query: json['query']! as String,
    sources: [
      for (final item in json['sources'] as List<Object?>? ?? const [])
        if (item case {'title': final String title}) title,
    ],
    createdAt: DateTime.parse(json['createdAt']! as String),
  );

  final String id;
  final String query;
  final List<String> sources;
  final DateTime createdAt;
}
