import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import 'deadline.dart';
import 'study_labels_context.dart';
import '../../design/components/tile_group.dart';

class DeadlineTile extends StatelessWidget {
  const DeadlineTile({required this.deadline, super.key});

  final Deadline deadline;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final submitted = deadline.submittedAt != null;
    return ListTile(
      contentPadding: TileGroup.tilePadding,
      leading: Icon(submitted ? Icons.task_alt : Icons.assignment_outlined),
      title: Text(deadline.title),
      subtitle: Text(
        [
          deadline.courseTitle,
          l10n.deadlineAt(labels.clock(deadline.dueAt)),
          if (submitted) l10n.assignmentSubmitted,
        ].join(', '),
      ),
      onTap: () => context.push(AppRoutes.assignment(deadline.assignmentId)),
    );
  }
}
