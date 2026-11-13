import '../../l10n/generated/app_localizations.dart';
import 'api_exception.dart';

abstract final class ApiErrorText {
  static String of(AppLocalizations l10n, Object error) {
    if (error is! ApiException) {
      return l10n.errorGeneric;
    }
    if (error.code == 'unknown' && error.statusCode == 404) {
      return l10n.errorNotLodestar;
    }
    return ofCode(l10n, error.code);
  }

  static String ofCode(AppLocalizations l10n, String code) {
    return switch (code) {
      ApiException.unreachable => l10n.errorUnreachable,
      ApiException.loginFailed || 'invalid_grant' => l10n.errorLoginFailed,
      'tenant_not_found' => l10n.errorTenantNotFound,
      'redirect_not_allowed' => l10n.errorRedirectRejected,
      'invalid_profile' => l10n.errorInvalidProfile,
      'deletion_by_university' => l10n.errorDeletionByUniversity,
      'notification_required' => l10n.channelRequired,
      'forbidden' => l10n.errorForbidden,
      'role_hierarchy' => l10n.errorRoleHierarchy,
      'permission_not_held' => l10n.errorPermissionNotHeld,
      'invalid_role' => l10n.errorInvalidRole,
      'unit_code_taken' => l10n.errorUnitCodeTaken,
      'unit_cycle' => l10n.errorUnitCycle,
      'unit_archived' => l10n.errorUnitArchived,
      'invalid_unit' => l10n.errorInvalidUnit,
      'invalid_assignment' => l10n.errorInvalidAssignment,
      'assignment_has_work' => l10n.errorAssignmentHasWork,
      'invalid_submission' => l10n.errorInvalidSubmission,
      'course_archived' => l10n.errorCourseArchived,
      'assignment_not_found' => l10n.errorAssignmentNotFound,
      'course_not_found' => l10n.errorCourseNotFound,
      'assignment_not_published' => l10n.errorNotPublished,
      'not_course_student' => l10n.errorNotStudent,
      'invalid_grade' => l10n.errorInvalidGrade,
      'note_not_found' => l10n.errorNoteNotFound,
      'invalid_note' => l10n.errorInvalidNote,
      'note_read_only' => l10n.errorNoteReadOnly,
      'note_deleted' => l10n.errorNoteDeleted,
      'notes_quota_exceeded' => l10n.errorNotesQuota,
      'not_note_owner' => l10n.errorNotNoteOwner,
      'stale_snapshot' => l10n.errorNoNewEdits,
      'invalid_source' => l10n.errorInvalidSource,
      'source_not_found' => l10n.errorSourceNotFound,
      'source_lookup_failed' => l10n.errorSourceLookup,
      'invalid_search_map' => l10n.errorInvalidSearchMap,
      'search_map_incomplete' => l10n.errorMapIncomplete,
      'navigator_unavailable' => l10n.errorNavigatorUnavailable,
      'navigator_failed' => l10n.errorNavigatorFailed,
      'dialog_not_found' => l10n.errorDialogNotFound,
      'scope_not_found' => l10n.errorScopeNotFound,
      'invalid_message' => l10n.errorInvalidMessage,
      'invalid_material' => l10n.errorInvalidMaterial,
      'material_not_found' => l10n.errorMaterialNotFound,
      'team_scope_not_found' || 'team_not_found' => l10n.errorTeamNotFound,
      'invalid_team' => l10n.errorInvalidTeam,
      'already_in_team' => l10n.errorAlreadyInTeam,
      'team_full' => l10n.errorTeamFull,
      'teams_by_teacher' => l10n.errorTeamsByTeacher,
      'works_alone' => l10n.errorWorksAlone,
      'teams_archived' => l10n.errorTeamsArchived,
      'not_invited' => l10n.errorNotInvited,
      'leave_request_pending' => l10n.errorLeavePending,
      'team_work' => l10n.errorTeamWork,
      'team_ai_unavailable' || 'team_ai_failed' => l10n.errorTeamAi,
      _ => l10n.errorGeneric,
    };
  }
}
