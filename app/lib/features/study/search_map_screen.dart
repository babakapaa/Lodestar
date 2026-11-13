import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'search_map_editor.dart';
import 'study_providers.dart';

class SearchMapScreen extends ConsumerWidget {
  const SearchMapScreen({required this.assignmentId, super.key});

  final String assignmentId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final draft = ref.watch(searchMapProvider(assignmentId));
    final error = draft.error;
    final loaded = draft.value?.value;
    final title = ref
        .watch(assignmentProvider(assignmentId))
        .value
        ?.value
        .assignment
        .title;
    if (loaded == null) {
      return Scaffold(
        appBar: AppBar(title: Text(l10n.mapTitle)),
        body: error is ApiException
            ? Padding(
                padding: const EdgeInsets.all(Spacing.l),
                child: FormError(message: ApiErrorText.of(l10n, error)),
              )
            : const Center(child: CircularProgressIndicator()),
      );
    }
    return SearchMapEditor(
      assignmentId: assignmentId,
      title: title ?? '',
      draft: loaded,
    );
  }
}
