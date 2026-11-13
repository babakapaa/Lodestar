class BlockedUser {
  const BlockedUser({required this.userId, required this.displayName});

  factory BlockedUser.fromJson(Map<String, Object?> json) => BlockedUser(
    userId: json['userId']! as String,
    displayName: json['displayName'] as String? ?? '',
  );

  final String userId;
  final String displayName;
}
