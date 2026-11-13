class LtiTool {
  const LtiTool({
    required this.loginUrl,
    required this.launchUrl,
    required this.jwksUrl,
  });

  factory LtiTool.fromJson(Map<String, Object?> json) => LtiTool(
    loginUrl: json['loginUrl']! as String,
    launchUrl: json['launchUrl']! as String,
    jwksUrl: json['jwksUrl']! as String,
  );

  final String loginUrl;
  final String launchUrl;
  final String jwksUrl;
}
