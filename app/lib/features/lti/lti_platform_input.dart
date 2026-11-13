class LtiPlatformInput {
  const LtiPlatformInput({
    required this.name,
    required this.issuer,
    required this.clientId,
    required this.authUrl,
    required this.tokenUrl,
    required this.jwksUrl,
    required this.deploymentIds,
    required this.matchEmail,
  });

  final String name;
  final String issuer;
  final String clientId;
  final String authUrl;
  final String tokenUrl;
  final String jwksUrl;
  final List<String> deploymentIds;
  final bool matchEmail;

  Map<String, Object?> toJson() => {
    'name': name,
    'issuer': issuer,
    'clientId': clientId,
    'authUrl': authUrl,
    'tokenUrl': tokenUrl,
    'jwksUrl': jwksUrl,
    'deploymentIds': deploymentIds,
    'matchEmail': matchEmail,
  };
}
