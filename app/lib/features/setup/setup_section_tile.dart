import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import '../people/people_providers.dart';
import 'setup_labels.dart';
import 'setup_section.dart';

class SetupSectionTile extends ConsumerWidget {
  const SetupSectionTile({required this.section, this.onTap, super.key});

  final SetupSection section;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final labels = SetupLabels(l10n);
    final assignee = section.assigneeId;
    final name = assignee == null
        ? null
        : ref.watch(personNameProvider(assignee)).value;
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Icon(
        section.isApproved ? Icons.check_circle : Icons.radio_button_unchecked,
        color: section.isApproved
            ? Theme.of(context).colorScheme.primary
            : null,
      ),
      title: Text(labels.section(section.key)),
      subtitle: Text(
        [
          labels.status(section.status),
          if (section.required) l10n.setupRequired,
          if (name != null) l10n.setupAssignee(name),
        ].join(', '),
      ),
      trailing: const Icon(Icons.chevron_right),
      onTap: onTap,
    );
  }
}
