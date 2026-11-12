class UnitMember {
  const UnitMember({
    required this.userId,
    required this.displayName,
    required this.email,
  });

  factory UnitMember.fromJson(Map<String, Object?> json) {
    return UnitMember(
      userId: json['userId']! as String,
      displayName: json['displayName']! as String,
      email: json['email']! as String,
    );
  }

  final String userId;
  final String displayName;
  final String email;
}
