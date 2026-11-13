class SandboxTemplate {
  const SandboxTemplate({
    required this.id,
    required this.name,
    required this.about,
    required this.image,
    required this.archived,
    required this.starter,
  });

  factory SandboxTemplate.fromJson(Map<String, Object?> json) =>
      SandboxTemplate(
        id: json['id']! as String,
        name: json['name'] as String? ?? '',
        about: json['about'] as String? ?? '',
        image: json['image'] as String? ?? '',
        archived: json['archived'] as bool? ?? false,
        starter: json['starter'] as bool? ?? false,
      );

  final String id;
  final String name;
  final String about;
  final String image;
  final bool archived;
  final bool starter;

  Map<String, Object?> toJson() => {
    'name': name,
    'about': about,
    'image': image,
    'archived': archived,
    'starter': starter,
  };
}
