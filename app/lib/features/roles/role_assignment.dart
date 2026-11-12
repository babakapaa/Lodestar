class RoleAssignment {
  const RoleAssignment({
    required this.id,
    required this.roleId,
    required this.scope,
    required this.expiresAt,
  });

  factory RoleAssignment.fromJson(Map<String, Object?> json) {
    final expiresAt = json['expiresAt'];
    return RoleAssignment(
      id: json['id']! as String,
      roleId: json['roleId']! as String,
      scope: json['scope']! as String,
      expiresAt: expiresAt is String ? DateTime.parse(expiresAt) : null,
    );
  }

  final String id;
  final String roleId;
  final String scope;
  final DateTime? expiresAt;

  bool get wholeUniversity => scope.isEmpty || scope == 'tenant';
}
