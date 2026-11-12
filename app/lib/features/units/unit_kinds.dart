import '../../l10n/generated/app_localizations.dart';

abstract final class UnitKinds {
  static const presets = ['faculty', 'department', 'stream', 'group'];

  static String label(AppLocalizations l10n, String kind) => switch (kind) {
    'faculty' => l10n.unitFaculty,
    'department' => l10n.unitDepartment,
    'stream' => l10n.unitStream,
    'group' => l10n.unitGroup,
    _ => kind,
  };
}
