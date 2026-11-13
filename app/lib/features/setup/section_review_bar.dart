import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'return_note_dialog.dart';
import 'setup_providers.dart';
import 'setup_section.dart';

class SectionReviewBar extends ConsumerWidget {
  const SectionReviewBar({
    required this.section,
    required this.manages,
    super.key,
  });

  final SetupSection section;
  final bool manages;

  Future<void> _run(
    BuildContext context,
    WidgetRef ref,
    Future<void> Function() action,
  ) async {
    await TeamAction.run(context, action);
    ref
      ..invalidate(setupOverviewProvider)
      ..invalidate(setupStatusProvider);
  }

  Future<void> _return(BuildContext context, WidgetRef ref) async {
    final note = await ReturnNoteDialog.show(context);
    if (note == null || !context.mounted) {
      return;
    }
    await _run(
      context,
      ref,
      () => ref
          .read(setupApiProvider)
          .review(section.key, approve: false, note: note),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final api = ref.read(setupApiProvider);
    return Wrap(
      spacing: Spacing.m,
      runSpacing: Spacing.m,
      children: [
        if (section.status != SetupSection.review && !section.isApproved)
          LodestarButton(
            key: const ValueKey('setup-submit'),
            label: l10n.setupSubmit,
            icon: Icons.send_outlined,
            variant: manages
                ? LodestarButtonVariant.secondary
                : LodestarButtonVariant.primary,
            onPressed: () => _run(context, ref, () => api.submit(section.key)),
          ),
        if (manages && !section.isApproved)
          LodestarButton(
            key: const ValueKey('setup-approve'),
            label: l10n.setupApprove,
            icon: Icons.check,
            onPressed: () => _run(
              context,
              ref,
              () => api.review(section.key, approve: true),
            ),
          ),
        if (manages && section.status != SetupSection.notStarted)
          LodestarButton(
            key: const ValueKey('setup-return'),
            label: l10n.setupReturn,
            icon: Icons.undo,
            variant: LodestarButtonVariant.secondary,
            onPressed: () => _return(context, ref),
          ),
      ],
    );
  }
}
