class RoleInput {
  const RoleInput({
    required this.name,
    required this.color,
    required this.permissions,
  });

  final String name;
  final String color;
  final Set<String> permissions;

  Map<String, Object?> toJson() => {
    'name': name,
    'color': color,
    'permissions': permissions.toList()..sort(),
  };
}
