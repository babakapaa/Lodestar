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
  static const chatReports = '/admin/chat-reports';
  static const awards = '/admin/awards';
  static const connect = '/connect';
  static const login = '/login';
  static const ltiLaunch = '/lti';
  static const ltiPlatforms = '/admin/lms';
  static const schedulePath = '/study/schedule';
  static const todayPath = '/today';
  static const notesPath = '/notes';
  static const chatsPath = '/chats';
  static const courseChatsPath = '/chats/course';
  static const notesTrash = '/notes/trash';
  static const notesGraph = '/notes/graph';
  static const sourcesPath = '/notes/sources';
  static const sourceNew = '/notes/sources/new';
  static const dialogNew = '/navigator/new';
  static const journal = '/navigator/journal';
  static const teams = '/teams';
  static const teamProfile = '/teams/profile';
  static const sandbox = '/sandbox';
  static const sandboxNew = '/sandbox/new';
  static const sandboxTopups = '/sandbox/topups';
  static const sandboxAdmin = '/admin/sandbox';
  static const attend = '/attend';
  static const stats = '/admin/stats';
  static const setup = '/admin/setup';

  static String role(String id) => '/admin/roles/$id';

  static String setupSection(String key) => '/admin/setup/$key';

  static String unitStats(String id, String name) => Uri(
    path: stats,
    queryParameters: {'parent': id, 'name': name},
  ).toString();

  static String ltiContext(String id) => '/lti/contexts/$id';

  static String courseLms(String id) => '/study/courses/$id/lms';

  static String courseSchedule(String id) => '/study/courses/$id/schedule';
  static String coursePreview(String id) => '/study/courses/$id/preview';

  static String personRoles(String id) => '/admin/people/$id';

  static String unit(String id) => '/admin/units/$id';

  static String course(String id) => '/study/courses/$id';

  static String assignmentNew(String courseId) =>
      '/study/courses/$courseId/new';

  static String assignment(String id) => '/study/assignments/$id';

  static String assignmentEdit(String id) => '/study/assignments/$id/edit';

  static String searchMap(String id) => '/study/assignments/$id/map';

  static String defenseBooking(String id) =>
      '/study/assignments/$id/defense/booking';

  static String defenses(String id) => '/study/assignments/$id/defenses';

  static String defense(String assignmentId, String studentId) =>
      '/study/assignments/$assignmentId/defenses/$studentId';

  static String grade(String assignmentId, String studentId) =>
      '/study/assignments/$assignmentId/students/$studentId';

  static String team(String id) => '/teams/$id';

  static String sandboxRun(String id) => '/sandbox/runs/$id';

  static String attendanceSession(String id) => '/attendance/sessions/$id';

  static String attendanceCourse(String id) => '/attendance/courses/$id';

  static String courseInsights(String id) => '/study/courses/$id/insights';

  static String sandboxNewFor(String kind, String id, String title) => Uri(
    path: sandboxNew,
    queryParameters: {'workKind': kind, 'workId': id, 'workTitle': title},
  ).toString();

  static String teamJoin(String code) => '/teams/codes/$code';

  static String teamScope(String kind, String id) => '/teams/scopes/$kind/$id';

  static String league(String kind, String id) => '/league/$kind/$id';

  static String leagueRules(String kind, String id) =>
      '/league/$kind/$id/rules';

  static String leaguePlayer(String kind, String id, String userId) =>
      '/league/$kind/$id/players/$userId';

  static String leagueTaskNew(String kind, String id) =>
      '/league/$kind/$id/tasks/new';

  static String leagueTask(String id) => '/league/tasks/$id';

  static String leagueTaskEdit(String id) => '/league/tasks/$id/edit';

  static String note(String id) => '/notes/$id';

  static String source(String id) => '/notes/sources/$id';

  static String noteMembers(String id) => '/notes/$id/members';

  static String noteVersions(String id) => '/notes/$id/versions';

  static String noteVersion(String id, String version) =>
      '/notes/$id/versions/$version';

  static String schedule(String date) => '$schedulePath?date=$date';

  static String dialog(String id) => '/navigator/dialogs/$id';

  static String chat(String id) => '/chats/$id';

  static String chatInfo(String id) => '/chats/$id/info';

  static String chatCall(String id) => '/chats/$id/call';

  static String courseChats(String courseId, {bool teach = false}) => Uri(
    path: '$courseChatsPath/$courseId',
    queryParameters: teach ? {'teach': '1'} : null,
  ).toString();

  static String newDialog({String? scopeKind, String? scopeId, String? text}) {
    final query = {
      if (scopeKind != null && scopeId != null) ...{
        'scopeKind': scopeKind,
        'scopeId': scopeId,
      },
      'text': ?text,
    };
    return Uri(
      path: dialogNew,
      queryParameters: query.isEmpty ? null : query,
    ).toString();
  }

  static String courseMaterials(String courseId) =>
      '/navigator/courses/$courseId/materials';

  static String materialNew(String courseId) =>
      '/navigator/courses/$courseId/materials/new';
}
