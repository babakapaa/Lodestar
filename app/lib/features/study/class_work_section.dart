import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment.dart';
import 'student_work_tile.dart';
import 'study_providers.dart';

class ClassWorkSection extends ConsumerWidget {
  const ClassWorkSection({required this.assignment, super.key});

  final Assignment assignment;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final work = ref.watch(assignmentWorkProvider(assignment.id));
    final items = work.value;
    final error = work.error;
    return Section(
      title: l10n.classWork,
      child: items != null
          ? items.isEmpty
                ? Text(l10n.classWorkNone)
                : Column(
                    children: [
                      for (final item in items)
                        StudentWorkTile(
                          key: ValueKey(item.student.userId),
                          work: item,
                          assignment: assignment,
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
