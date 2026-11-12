class PersonHit {
  const PersonHit({
    required this.id,
    required this.name,
    required this.subtitle,
  });

  factory PersonHit.fromJson(Map<String, Object?> json) {
    return PersonHit(
      id: json['id']! as String,
      name: json['title']! as String,
      subtitle: json['subtitle']! as String,
    );
  }

  final String id;
  final String name;
  final String subtitle;
}
