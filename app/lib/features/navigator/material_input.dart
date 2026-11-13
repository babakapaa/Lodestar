class MaterialInput {
  const MaterialInput({
    required this.title,
    required this.reference,
    required this.url,
    required this.text,
  });

  final String title;
  final String reference;
  final String url;
  final String text;

  Map<String, Object?> toJson() => {
    'title': title,
    'reference': reference,
    'url': url,
    'text': text,
  };
}
