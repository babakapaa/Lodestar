class PermissionSet {
  const PermissionSet(this._granted);

  static const empty = PermissionSet({});
  static const administrator = 'core.administrator';
  static const manageSettings = 'core.settings.manage';
  static const manageRoles = 'core.roles.manage';
  static const manageStructure = 'core.structure.manage';
  static const manageMembers = 'core.structure.members';
  static const viewAudit = 'core.audit.view';
  static const reviewChatReports = 'chat.reports.review';
  static const manageAwards = 'league.awards.manage';
  static const manageIntegrations = 'core.integrations.manage';
  static const manageSandbox = 'sandbox.manage';
  static const viewStats = 'core.stats.view';
  static const manageSetup = 'core.setup.manage';

  final Set<String> _granted;

  bool has(String permission) =>
      _granted.contains(administrator) || _granted.contains(permission);

  bool get administers =>
      has(manageSettings) ||
      has(viewAudit) ||
      has(manageRoles) ||
      has(manageStructure) ||
      has(manageMembers) ||
      has(reviewChatReports) ||
      has(manageAwards) ||
      has(manageIntegrations) ||
      has(manageSandbox) ||
      has(viewStats) ||
      has(manageSetup);
}
