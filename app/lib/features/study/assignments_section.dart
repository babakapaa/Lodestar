import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_tile.dart';
import 'study_providers.dart';

class AssignmentsSection extends ConsumerWidget {
  const AssignmentsSection({
    required this.courseId,
    required this.teaches,
    super.key,
  });

  final String courseId;
  final bool teaches;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final assignments = ref.watch(assignmentsProvider(courseId));
    final items = assignments.value?.value;
    final error = assignments.error;
    return Section(
      title: l10n.courseAssignments,
      child: items != null
          ? items.isEmpty
                ? Text(l10n.courseNoAssignments)
                : Column(
                    children: [
                      for (final item in items)
                        AssignmentTile(
                          key: ValueKey(item.assignment.id),
                          item: item,
                          teaches: teaches,
                        ),
                    ],
                  )
          : error != null
          ? FormError(message: ApiErrorText.of(l10n, error))
          : const Padding(
              padding: EdgeInsets.all(Spacing.l),
              child: Center(child: CircularProgressIndicator()),
            ),
    );
  }
}
