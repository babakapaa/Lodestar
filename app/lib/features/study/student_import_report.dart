import 'package:material_ui/material_ui.dart';

import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/generated/app_localizations.dart';
import '../../l10n/l10n_context.dart';
import 'student_import.dart';
import 'student_import_error.dart';

class StudentImportReport extends StatelessWidget {
  const StudentImportReport({required this.report, super.key});

  final StudentImport report;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    if (report.errors.isEmpty) {
      return Text(
        key: const ValueKey('import-done'),
        l10n.importDone(report.added, report.existing),
      );
    }
    return Column(
      key: const ValueKey('import-errors'),
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        NoticeBanner(
          icon: Icons.error_outline,
          message: l10n.importRejected(report.errors.length),
        ),
        const SizedBox(height: Spacing.s),
        for (final error in report.errors)
          Text(l10n.importLine(error.line, _reason(l10n, error))),
      ],
    );
  }

  String _reason(AppLocalizations l10n, StudentImportError error) =>
      switch (error.code) {
        'missing_email' => l10n.importMissingEmail,
        'unknown_user' => l10n.importUnknownUser,
        'ambiguous_user' => l10n.importAmbiguousUser,
        _ => error.code,
      };
}
