class UnitsQuery {
  const UnitsQuery({required this.parentId, required this.archived});

  final String? parentId;
  final bool archived;

  @override
  bool operator ==(Object other) =>
      other is UnitsQuery &&
      other.parentId == parentId &&
      other.archived == archived;

  @override
  int get hashCode => Object.hash(parentId, archived);
}
