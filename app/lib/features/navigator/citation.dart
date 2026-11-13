class Citation {
  const Citation({
    required this.n,
    required this.title,
    required this.reference,
    required this.url,
    required this.text,
  });

  factory Citation.fromJson(Map<String, Object?> json) => Citation(
    n: json['n']! as int,
    title: json['title']! as String,
    reference: json['reference'] as String? ?? '',
    url: json['url'] as String? ?? '',
    text: json['text'] as String? ?? '',
  );

  final int n;
  final String title;
  final String reference;
  final String url;
  final String text;
}
