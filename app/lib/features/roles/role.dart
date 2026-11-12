class Role {
  const Role({
    required this.id,
    required this.name,
    required this.color,
    required this.position,
    required this.permissions,
    required this.isDefault,
  });

  factory Role.fromJson(Map<String, Object?> json) {
    return Role(
      id: json['id']! as String,
      name: json['name']! as String,
      color: json['color']! as String,
      position: (json['position']! as num).toInt(),
      permissions: {
        for (final item in json['permissions']! as List<Object?>)
          item! as String,
      },
      isDefault: json['default']! as bool,
    );
  }

  final String id;
  final String name;
  final String color;
  final int position;
  final Set<String> permissions;
  final bool isDefault;
}
