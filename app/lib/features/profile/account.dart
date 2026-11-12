class Account {
  const Account({required this.email, required this.displayName});

  factory Account.fromJson(Map<String, Object?> json) {
    return Account(
      email: json['email']! as String,
      displayName: json['displayName']! as String,
    );
  }

  final String email;
  final String displayName;
}
