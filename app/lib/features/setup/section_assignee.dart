import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import '../people/people_providers.dart';
import '../people/person_picker_dialog.dart';
import '../team/team_action.dart';
import 'setup_providers.dart';
import 'setup_section.dart';

class SectionAssignee extends ConsumerWidget {
  const SectionAssignee({required this.section, super.key});

  final SetupSection section;

  Future<void> _assign(
    BuildContext context,
    WidgetRef ref,
    String? userId,
  ) async {
    await TeamAction.run(
      context,
      () => ref.read(setupApiProvider).assign(section.key, userId),
    );
    ref.invalidate(setupOverviewProvider);
  }

  Future<void> _pick(BuildContext context, WidgetRef ref) async {
    final person = await PersonPickerDialog.show(context);
    if (person == null || !context.mounted) {
      return;
    }
    await _assign(context, ref, person.id);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final assignee = section.assigneeId;
    final name = assignee == null
        ? null
        : ref.watch(personNameProvider(assignee)).value;
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: const Icon(Icons.assignment_ind_outlined),
      title: Text(name == null ? l10n.setupAssign : l10n.setupAssignee(name)),
      trailing: assignee == null
          ? const Icon(Icons.chevron_right)
          : IconButton(
              tooltip: l10n.setupUnassign,
              icon: const Icon(Icons.close),
              onPressed: () => _assign(context, ref, null),
            ),
      onTap: () => _pick(context, ref),
    );
  }
}
