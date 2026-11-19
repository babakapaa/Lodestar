class Article {
  const Article({
    required this.title,
    required this.authors,
    required this.year,
    required this.venue,
    required this.doi,
    required this.url,
  });

  factory Article.fromJson(Map<String, Object?> json) => Article(
    title: json['title'] as String? ?? '',
    authors: [
      for (final item in json['authors'] as List<Object?>? ?? const [])
        item! as String,
    ],
    year: json['year'] as int? ?? 0,
    venue: json['venue'] as String? ?? '',
    doi: json['doi'] as String? ?? '',
    url: json['url'] as String? ?? '',
  );

  final String title;
  final List<String> authors;
  final int year;
  final String venue;
  final String doi;
  final String url;

  String get byline => [
    if (authors.isNotEmpty) authors.take(3).join(', '),
    if (year > 0) '$year',
    if (venue.isNotEmpty) venue,
  ].join('. ');
}
