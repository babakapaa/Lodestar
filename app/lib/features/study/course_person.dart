class CoursePerson {
  const CoursePerson({
    required this.userId,
    required this.displayName,
    required this.email,
  });

  factory CoursePerson.fromJson(Map<String, Object?> json) {
    return CoursePerson(
      userId: json['userId']! as String,
      displayName: json['displayName']! as String,
      email: json['email']! as String,
    );
  }

  final String userId;
  final String displayName;
  final String email;

  String get label => displayName.isEmpty ? email : displayName;
}
