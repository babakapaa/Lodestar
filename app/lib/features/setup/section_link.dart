import '../../core/navigation/app_routes.dart';

abstract final class SectionLink {
  static String? of(String section) => switch (section) {
    'university' => AppRoutes.units,
    'people' || 'admins' => AppRoutes.roles,
    'integrations' => AppRoutes.ltiPlatforms,
    'league' => AppRoutes.awards,
    'sandbox' => AppRoutes.sandboxAdmin,
    _ => null,
  };
}
