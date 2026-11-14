import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../notes/note_actions.dart';

class TodayActions extends ConsumerWidget {
  const TodayActions({required this.date, super.key});

  final String date;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Wrap(
      spacing: Spacing.s,
      runSpacing: Spacing.s,
      children: [
        OutlinedButton.icon(
          onPressed: () => context.push(AppRoutes.newDialog()),
          icon: const Icon(Icons.explore_outlined),
          label: Text(l10n.todayAsk),
        ),
        OutlinedButton.icon(
          onPressed: () => NoteActions.create(context, ref),
          icon: const Icon(Icons.edit_note),
          label: Text(l10n.notesNew),
        ),
        OutlinedButton.icon(
          onPressed: () => context.push(AppRoutes.schedule(date)),
          icon: const Icon(Icons.calendar_month_outlined),
          label: Text(l10n.todayWeek),
        ),
      ],
    );
  }
}
