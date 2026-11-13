import 'link_target.dart';

class LtiDeepLink {
  const LtiDeepLink({
    required this.id,
    required this.platform,
    required this.acceptMultiple,
    required this.targets,
  });

  factory LtiDeepLink.fromJson(Map<String, Object?> json) => LtiDeepLink(
    id: json['id']! as String,
    platform: json['platform'] as String? ?? '',
    acceptMultiple: json['acceptMultiple'] as bool? ?? false,
    targets: [
      for (final item in json['targets'] as List<Object?>? ?? const [])
        LinkTarget.fromJson(item! as Map<String, Object?>),
    ],
  );

  final String id;
  final String platform;
  final bool acceptMultiple;
  final List<LinkTarget> targets;
}
