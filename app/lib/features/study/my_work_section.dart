import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/notice_banner.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_details.dart';
import 'study_labels_context.dart';
import 'study_providers.dart';

class MyWorkSection extends ConsumerWidget {
  const MyWorkSection({required this.details, super.key});

  final AssignmentDetails details;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final textTheme = Theme.of(context).textTheme;
    final muted = textTheme.bodyMedium!.copyWith(color: context.palette.muted);
    final waiting = ref.watch(
      outboxProvider.select(
        (drafts) => drafts.any((d) => d.assignmentId == details.assignment.id),
      ),
    );
    final mine = details.mine;
    final latest = mine?.latest;
    final grade = mine?.grade;
    return Section(
      title: l10n.myWork,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (waiting) ...[
            NoticeBanner(
              icon: Icons.schedule_send_outlined,
              message: l10n.submitQueued,
            ),
            const SizedBox(height: Spacing.m),
          ],
          if (latest == null)
            Text(l10n.myWorkNone, style: muted)
          else
            LodestarCard(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    [
                      l10n.myWorkSentAt(labels.time(latest.submittedAt)),
                      if (mine!.attempts > 1)
                        l10n.myWorkAttempts(mine.attempts),
                    ].join(', '),
                    style: muted,
                  ),
                  if (latest.body.isNotEmpty) ...[
                    const SizedBox(height: Spacing.s),
                    SelectableText(latest.body),
                  ],
                  for (final file in latest.files)
                    ListTile(
                      contentPadding: EdgeInsets.zero,
                      leading: const Icon(Icons.attach_file),
                      title: Text(file.name),
                    ),
                ],
              ),
            ),
          if (grade != null) ...[
            const SizedBox(height: Spacing.m),
            LodestarCard(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    labels.grade(grade, details.assignment),
                    style: textTheme.titleMedium,
                  ),
                  if (grade.comment.isNotEmpty) ...[
                    const SizedBox(height: Spacing.s),
                    Text(l10n.feedback, style: muted),
                    const SizedBox(height: Spacing.xs),
                    SelectableText(grade.comment),
                  ],
                ],
              ),
            ),
          ],
        ],
      ),
    );
  }
}
