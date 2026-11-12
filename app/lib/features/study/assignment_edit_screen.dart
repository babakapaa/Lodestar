import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_form.dart';
import 'study_providers.dart';

class AssignmentEditScreen extends ConsumerWidget {
  const AssignmentEditScreen({this.courseId, this.assignmentId, super.key})
    : assert((courseId == null) != (assignmentId == null));

  final String? courseId;
  final String? assignmentId;

  Widget _content(BuildContext context, WidgetRef ref) {
    final courseId = this.courseId;
    if (courseId != null) {
      return AssignmentForm(courseId: courseId);
    }
    final details = ref.watch(assignmentProvider(assignmentId!));
    final loaded = details.value?.value;
    final error = details.error;
    if (loaded != null) {
      return AssignmentForm(
        key: ValueKey(loaded.assignment.id),
        courseId: loaded.assignment.courseId,
        details: loaded,
      );
    }
    if (error != null) {
      return Padding(
        padding: const EdgeInsets.all(Spacing.l),
        child: FormError(message: ApiErrorText.of(context.l10n, error)),
      );
    }
    return const Center(child: CircularProgressIndicator());
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Scaffold(
      appBar: AppBar(
        title: Text(
          assignmentId == null ? l10n.assignmentNew : l10n.assignmentEdit,
        ),
      ),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: _content(context, ref),
          ),
        ),
      ),
    );
  }
}
