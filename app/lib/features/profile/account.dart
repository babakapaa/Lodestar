class Account {
  const Account({
    required this.id,
    required this.email,
    required this.displayName,
  });

  factory Account.fromJson(Map<String, Object?> json) {
    return Account(
      id: json['id'] as String? ?? '',
      email: json['email']! as String,
      displayName: json['displayName']! as String,
    );
  }

  final String id;
  final String email;
  final String displayName;
}
