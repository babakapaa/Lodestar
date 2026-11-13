import 'lti_platform_input.dart';

class LtiPlatform {
  const LtiPlatform({required this.id, required this.input});

  factory LtiPlatform.fromJson(Map<String, Object?> json) => LtiPlatform(
    id: json['id']! as String,
    input: LtiPlatformInput(
      name: json['name'] as String? ?? '',
      issuer: json['issuer'] as String? ?? '',
      clientId: json['clientId'] as String? ?? '',
      authUrl: json['authUrl'] as String? ?? '',
      tokenUrl: json['tokenUrl'] as String? ?? '',
      jwksUrl: json['jwksUrl'] as String? ?? '',
      deploymentIds: [
        for (final id in json['deploymentIds'] as List<Object?>? ?? const [])
          id! as String,
      ],
      matchEmail: json['matchEmail'] as bool? ?? false,
    ),
  );

  final String id;
  final LtiPlatformInput input;
}
