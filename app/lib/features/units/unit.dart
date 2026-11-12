class Unit {
  const Unit({
    required this.id,
    required this.parentId,
    required this.kind,
    required this.name,
    required this.code,
    required this.archived,
  });

  factory Unit.fromJson(Map<String, Object?> json) {
    return Unit(
      id: json['id']! as String,
      parentId: json['parentId'] as String?,
      kind: json['kind']! as String,
      name: json['name']! as String,
      code: json['code'] as String?,
      archived: json['archivedAt'] != null,
    );
  }

  final String id;
  final String? parentId;
  final String kind;
  final String name;
  final String? code;
  final bool archived;
}
