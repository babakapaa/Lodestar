class NoteMember {
  const NoteMember({
    required this.userId,
    required this.displayName,
    required this.email,
    required this.canEdit,
  });

  factory NoteMember.fromJson(Map<String, Object?> json) => NoteMember(
    userId: json['userId']! as String,
    displayName: json['displayName']! as String,
    email: json['email']! as String,
    canEdit: json['canEdit']! as bool,
  );

  final String userId;
  final String displayName;
  final String email;
  final bool canEdit;

  String get label => displayName.isEmpty ? email : displayName;
}
