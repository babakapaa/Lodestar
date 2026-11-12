class IdentityProvider {
  const IdentityProvider({required this.key, required this.displayName});

  factory IdentityProvider.fromJson(Map<String, Object?> json) {
    return IdentityProvider(
      key: json['key']! as String,
      displayName: json['displayName']! as String,
    );
  }

  final String key;
  final String displayName;
}
