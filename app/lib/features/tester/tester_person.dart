class TesterPerson {
  const TesterPerson({
    required this.key,
    required this.name,
    required this.title,
    required this.email,
  });

  factory TesterPerson.fromJson(Map<String, Object?> json) {
    return TesterPerson(
      key: json['key']! as String,
      name: json['name']! as String,
      title: json['title']! as String,
      email: json['email']! as String,
    );
  }

  final String key;
  final String name;
  final String title;
  final String email;
}
