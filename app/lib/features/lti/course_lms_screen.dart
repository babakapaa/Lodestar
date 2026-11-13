import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../study/gradebook_section.dart';
import '../study/study_providers.dart';
import 'lti_context_card.dart';
import 'lti_providers.dart';

class CourseLmsScreen extends ConsumerWidget {
  const CourseLmsScreen({required this.courseId, super.key});

  final String courseId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    return Scaffold(
      appBar: AppBar(title: Text(l10n.courseTablesTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(courseLmsProvider(courseId)),
          onRefresh: () => ref.refresh(courseLmsProvider(courseId).future),
          builder: (contexts) => [
            GradebookSection(
              courseId: courseId,
              onImported: () => ref.invalidate(courseProvider(courseId)),
            ),
            const SizedBox(height: Spacing.xl),
            Section(
              title: l10n.lmsTitle,
              child: contexts.isEmpty
                  ? Text(l10n.lmsCourseNone, style: muted)
                  : Column(
                      children: [
                        for (final c in contexts) ...[
                          LtiContextCard(
                            context: c,
                            openLink: false,
                            onChanged: () =>
                                ref.invalidate(courseLmsProvider(courseId)),
                          ),
                          const SizedBox(height: Spacing.m),
                        ],
                      ],
                    ),
            ),
          ],
        ),
      ),
    );
  }
}
