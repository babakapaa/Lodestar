import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'assignment.dart';
import 'attached_file_tile.dart';
import 'grade_form.dart';
import 'help_summary_view.dart';
import 'search_map_view.dart';
import 'student_work.dart';
import 'study_labels_context.dart';
import 'study_providers.dart';
import 'submission.dart';

class GradeScreen extends ConsumerWidget {
  const GradeScreen({
    required this.assignmentId,
    required this.studentId,
    super.key,
  });

  final String assignmentId;
  final String studentId;

  List<Widget> _attempts(BuildContext context, List<Submission> attempts) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    if (attempts.isEmpty) {
      return [Text(l10n.gradeNoWork, style: muted)];
    }
    return [
      for (final attempt in attempts) ...[
        LodestarCard(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                l10n.myWorkSentAt(labels.time(attempt.submittedAt)),
                style: muted,
              ),
              if (attempt.body.isNotEmpty) ...[
                const SizedBox(height: Spacing.s),
                SelectableText(attempt.body),
              ],
              for (final file in attempt.files) AttachedFileTile(file: file),
              if (attempt.searchMap case final map?) ...[
                const SizedBox(height: Spacing.m),
                SearchMapView(map: map),
              ],
              if (attempt.help case final help?) ...[
                const SizedBox(height: Spacing.m),
                HelpSummaryView(summary: help),
              ],
            ],
          ),
        ),
        const SizedBox(height: Spacing.m),
      ],
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final details = ref.watch(assignmentProvider(assignmentId));
    final work = ref.watch(assignmentWorkProvider(assignmentId));
    final attempts = ref.watch(attemptsProvider((assignmentId, studentId)));
    final error = details.error ?? work.error ?? attempts.error;
    final Assignment? assignment = details.value?.value.assignment;
    StudentWork? student;
    for (final item in work.value ?? const <StudentWork>[]) {
      if (item.student.userId == studentId) {
        student = item;
      }
    }
    final list = attempts.value;
    return Scaffold(
      appBar: AppBar(title: Text(student?.student.label ?? '')),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 720),
            child: error != null
                ? Padding(
                    padding: const EdgeInsets.all(Spacing.l),
                    child: FormError(message: ApiErrorText.of(l10n, error)),
                  )
                : assignment == null || student == null || list == null
                ? const Center(child: CircularProgressIndicator())
                : ListView(
                    padding: const EdgeInsets.all(Spacing.l),
                    children: [
                      Text(
                        assignment.title,
                        style: Theme.of(context).textTheme.titleLarge,
                      ),
                      const SizedBox(height: Spacing.xl),
                      Section(
                        title: l10n.gradeAttempts,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.stretch,
                          children: _attempts(context, list),
                        ),
                      ),
                      const SizedBox(height: Spacing.xl),
                      Section(
                        title: l10n.gradeTitle,
                        child: GradeForm(
                          assignment: assignment,
                          studentId: studentId,
                          grade: student.grade,
                        ),
                      ),
                    ],
                  ),
          ),
        ),
      ),
    );
  }
}
