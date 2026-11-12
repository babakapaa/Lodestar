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
  static const schedulePath = '/study/schedule';
  static const notesPath = '/notes';
  static const notesTrash = '/notes/trash';

  static String role(String id) => '/admin/roles/$id';

  static String personRoles(String id) => '/admin/people/$id';

  static String unit(String id) => '/admin/units/$id';

  static String course(String id) => '/study/courses/$id';

  static String assignmentNew(String courseId) =>
      '/study/courses/$courseId/new';

  static String assignment(String id) => '/study/assignments/$id';

  static String assignmentEdit(String id) => '/study/assignments/$id/edit';

  static String grade(String assignmentId, String studentId) =>
      '/study/assignments/$assignmentId/students/$studentId';

  static String note(String id) => '/notes/$id';

  static String noteMembers(String id) => '/notes/$id/members';

  static String schedule(String date) => '$schedulePath?date=$date';
}
