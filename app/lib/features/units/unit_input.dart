class UnitInput {
  const UnitInput({
    required this.kind,
    required this.name,
    required this.code,
    this.parentId,
  });

  final String? parentId;
  final String kind;
  final String name;
  final String code;

  Map<String, Object?> toJson() => {
    'parentId': ?parentId,
    'kind': kind,
    'name': name,
    'code': code.isEmpty ? null : code,
  };
}
