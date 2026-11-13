import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'setup_providers.dart';

class PilotCard extends ConsumerWidget {
  const PilotCard({required this.courses, super.key});

  final int courses;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Section(
      title: l10n.setupPilotTitle,
      child: ListTile(
        contentPadding: EdgeInsets.zero,
        leading: const Icon(Icons.science_outlined),
        title: Text(l10n.setupPilotCourses(courses)),
        trailing: TextButton(
          key: const ValueKey('setup-expand'),
          onPressed: () async {
            await TeamAction.run(context, ref.read(setupApiProvider).expand);
            ref
              ..invalidate(setupOverviewProvider)
              ..invalidate(setupStatusProvider);
          },
          child: Text(l10n.setupExpand),
        ),
      ),
    );
  }
}
