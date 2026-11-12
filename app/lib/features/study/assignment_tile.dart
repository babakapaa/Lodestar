import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_item.dart';
import 'study_labels_context.dart';
import 'study_providers.dart';

class AssignmentTile extends ConsumerWidget {
  const AssignmentTile({required this.item, required this.teaches, super.key});

  final AssignmentItem item;
  final bool teaches;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final assignment = item.assignment;
    final waiting = ref.watch(
      outboxProvider.select(
        (drafts) => drafts.any((d) => d.assignmentId == assignment.id),
      ),
    );
    final grade = item.grade;
    final counts = item.counts;
    final dueAt = assignment.dueAt;
    final overdue =
        !teaches &&
        item.submittedAt == null &&
        !waiting &&
        dueAt != null &&
        dueAt.isBefore(DateTime.now());
    final status = teaches
        ? [
            if (!assignment.published) l10n.assignmentDraft,
            if (counts != null)
              l10n.assignmentProgress(
                counts.submitted,
                counts.students,
                counts.graded,
              ),
          ]
        : [
            if (waiting)
              l10n.assignmentWaiting
            else if (grade != null)
              labels.grade(grade, assignment)
            else if (item.submittedAt != null)
              l10n.assignmentSubmitted
            else if (overdue)
              l10n.assignmentOverdue,
          ];
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Icon(
        waiting
            ? Icons.schedule_send_outlined
            : grade != null || item.submittedAt != null
            ? Icons.task_alt
            : Icons.assignment_outlined,
        color: overdue ? context.palette.warn : null,
      ),
      title: Text(assignment.title),
      subtitle: Text([labels.due(dueAt), ...status].join('\n')),
      isThreeLine: status.isNotEmpty,
      trailing: const Icon(Icons.chevron_right),
      onTap: () => context.push(AppRoutes.assignment(assignment.id)),
    );
  }
}
