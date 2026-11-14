import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';

class NoteNotices extends StatelessWidget {
  const NoteNotices({required this.note, super.key});

  final Note note;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = note.error;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        for (final message in [
          if (note.deleted) l10n.notesInTrash,
          if (!note.role.canEdit) l10n.notesReadOnly,
          if (error != null)
            l10n.notesNotAccepted(ApiErrorText.ofCode(l10n, error)),
        ])
          Padding(
            padding: const EdgeInsets.fromLTRB(
              Spacing.l,
              Spacing.m,
              Spacing.l,
              0,
            ),
            child: NoticeBanner(icon: Icons.info_outline, message: message),
          ),
      ],
    );
  }
}
