class TokenPair {
  const TokenPair({required this.accessToken, required this.refreshToken});

  factory TokenPair.fromJson(Map<String, Object?> json) {
    return TokenPair(
      accessToken: json['accessToken']! as String,
      refreshToken: json['refreshToken']! as String,
    );
  }

  final String accessToken;
  final String refreshToken;

  Map<String, Object?> toJson() => {
    'accessToken': accessToken,
    'refreshToken': refreshToken,
  };
}
