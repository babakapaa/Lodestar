import '../../l10n/generated/app_localizations.dart';
import 'api_exception.dart';

abstract final class ApiErrorText {
  static String of(AppLocalizations l10n, Object error) {
    if (error is! ApiException) {
      return l10n.errorGeneric;
    }
    return switch (error.code) {
      ApiException.unreachable => l10n.errorUnreachable,
      ApiException.loginFailed || 'invalid_grant' => l10n.errorLoginFailed,
      'tenant_not_found' => l10n.errorTenantNotFound,
      'redirect_not_allowed' => l10n.errorRedirectRejected,
      'invalid_profile' => l10n.errorInvalidProfile,
      'deletion_by_university' => l10n.errorDeletionByUniversity,
      'notification_required' => l10n.channelRequired,
      'unknown' when error.statusCode == 404 => l10n.errorNotLodestar,
      _ => l10n.errorGeneric,
    };
  }
}
