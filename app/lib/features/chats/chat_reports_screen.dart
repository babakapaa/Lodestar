import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'chat_remote.dart';
import 'chat_report_card.dart';
import 'chats_providers.dart';

class ChatReportsScreen extends ConsumerWidget {
  const ChatReportsScreen({super.key});

  Future<void> _resolve(BuildContext context, WidgetRef ref, String id) async {
    await ChatRemote.run(context, () => ref.read(chatsApiProvider).resolve(id));
    ref.invalidate(chatReportsProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final reports = ref.watch(chatReportsProvider);
    final list = reports.value;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.chatReportsTitle)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: RefreshIndicator(
            onRefresh: () => ref.refresh(chatReportsProvider.future),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: [
                if (reports.hasError)
                  FormError(message: ApiErrorText.of(l10n, reports.error!))
                else if (list == null)
                  const Center(child: CircularProgressIndicator())
                else if (list.isEmpty)
                  EmptyState(
                    icon: Icons.flag_outlined,
                    title: l10n.chatReportsEmptyTitle,
                    body: l10n.chatReportsEmptyBody,
                  )
                else
                  for (final report in list) ...[
                    ChatReportCard(
                      key: ValueKey(report.id),
                      report: report,
                      onResolve: () => _resolve(context, ref, report.id),
                    ),
                    const SizedBox(height: Spacing.m),
                  ],
              ],
            ),
          ),
        ),
      ),
    );
  }
}
