abstract final class AppRoutes {
  static const settings = '/settings';
  static const profile = '/profile';
  static const profileEdit = '/profile/edit';
  static const privacy = '/profile/privacy';
  static const notifications = '/notifications';
  static const notificationSettings = '/notifications/settings';
  static const admin = '/admin';
  static const adminSettings = '/admin/settings';
  static const adminAudit = '/admin/audit';
  static const roles = '/admin/roles';
  static const roleNew = '/admin/roles/new';
  static const people = '/admin/people';
  static const units = '/admin/units';
  static const connect = '/connect';
  static const login = '/login';

  static String role(String id) => '/admin/roles/$id';

  static String personRoles(String id) => '/admin/people/$id';

  static String unit(String id) => '/admin/units/$id';
}
