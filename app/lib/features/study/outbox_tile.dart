import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'study_providers.dart';
import 'submission_draft.dart';

class OutboxTile extends ConsumerWidget {
  const OutboxTile({required this.draft, super.key});

  final SubmissionDraft draft;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final error = draft.error;
    final outbox = ref.read(outboxProvider.notifier);
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Icon(
        error == null ? Icons.schedule_send_outlined : Icons.error_outline,
        color: error == null ? null : context.palette.warn,
      ),
      title: Text(draft.assignmentTitle),
      subtitle: Text(
        error == null
            ? l10n.studyOutboxPending
            : l10n.studyOutboxFailed(ApiErrorText.ofCode(l10n, error)),
      ),
      onTap: () => context.push(AppRoutes.assignment(draft.assignmentId)),
      trailing: error == null
          ? null
          : Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                IconButton(
                  tooltip: l10n.studyRetry,
                  icon: const Icon(Icons.refresh),
                  onPressed: () => outbox.retry(draft.clientId),
                ),
                IconButton(
                  tooltip: l10n.studyDiscard,
                  icon: const Icon(Icons.delete_outline),
                  onPressed: () => outbox.discard(draft.clientId),
                ),
              ],
            ),
    );
  }
}
