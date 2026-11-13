import '../../l10n/generated/app_localizations.dart';
import 'setup_problem.dart';
import 'setup_section.dart';

class SetupLabels {
  const SetupLabels(this.l10n);

  final AppLocalizations l10n;

  String section(String key) => switch (key) {
    'deploy' => l10n.setupDeploy,
    'university' => l10n.setupUniversity,
    'people' => l10n.setupPeople,
    'admins' => l10n.setupAdmins,
    'integrations' => l10n.setupIntegrations,
    'teaching' => l10n.setupTeaching,
    'navigator' => l10n.setupNavigator,
    'league' => l10n.setupLeague,
    'sandbox' => l10n.setupSandbox,
    'campus' => l10n.setupCampus,
    'privacy' => l10n.setupPrivacy,
    'legal' => l10n.setupLegal,
    'launch' => l10n.setupLaunch,
    _ => key,
  };

  String hint(String key) => switch (key) {
    'deploy' => l10n.setupDeployHint,
    'university' => l10n.setupUniversityHint,
    'people' => l10n.setupPeopleHint,
    'admins' => l10n.setupAdminsHint,
    'integrations' => l10n.setupIntegrationsHint,
    'teaching' => l10n.setupTeachingHint,
    'navigator' => l10n.setupNavigatorHint,
    'league' => l10n.setupLeagueHint,
    'sandbox' => l10n.setupSandboxHint,
    'campus' => l10n.setupCampusHint,
    'privacy' => l10n.setupPrivacyHint,
    'legal' => l10n.setupLegalHint,
    'launch' => l10n.setupLaunchHint,
    _ => '',
  };

  String status(String status) => switch (status) {
    SetupSection.inProgress => l10n.setupInProgress,
    SetupSection.review => l10n.setupInReview,
    SetupSection.approved => l10n.setupApproved,
    _ => l10n.setupNotStarted,
  };

  String question(String key, String fallback) => switch (key) {
    'storage' => l10n.setupQuestionStorage,
    'consent' => l10n.setupQuestionConsent,
    'red_button' => l10n.setupQuestionRedButton,
    'recording' => l10n.setupQuestionRecording,
    'retention' => l10n.setupQuestionRetention,
    'personal_access' => l10n.setupQuestionPersonalAccess,
    'sandbox_watch' => l10n.setupQuestionSandboxWatch,
    'data_officer' => l10n.setupQuestionDataOfficer,
    _ => fallback,
  };

  String problem(SetupProblem p) => switch (p.code) {
    'section_not_approved' => l10n.setupProblemNotApproved(section(p.section)),
    'legal_open' => l10n.setupProblemLegalOpen(p.count),
    'protocol_unsigned' => l10n.setupProblemUnsigned,
    'no_units' => l10n.setupProblemNoUnits,
    'courses_without_teacher' => l10n.setupProblemLonelyCourses(p.count),
    'model_missing' => l10n.setupProblemModel,
    _ => l10n.setupProblemOther(p.code),
  };
}
