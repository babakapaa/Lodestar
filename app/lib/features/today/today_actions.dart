import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../notes/note_actions.dart';
import 'today_action_card.dart';

class TodayActions extends ConsumerWidget {
  const TodayActions({required this.date, super.key});

  final String date;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final cards = [
      TodayActionCard(
        icon: Icons.explore_outlined,
        title: l10n.todayAsk,
        subtitle: l10n.todayAskHint,
        onTap: () => context.push(AppRoutes.newDialog()),
      ),
      TodayActionCard(
        icon: Icons.edit_note,
        title: l10n.notesNew,
        subtitle: l10n.todayNoteHint,
        onTap: () => NoteActions.create(context, ref),
      ),
      TodayActionCard(
        icon: Icons.calendar_month_outlined,
        title: l10n.todayWeek,
        subtitle: l10n.todayWeekHint,
        onTap: () => context.push(AppRoutes.schedule(date)),
      ),
    ];
    return LayoutBuilder(
      builder: (context, constraints) {
        if (constraints.maxWidth < 640) {
          return Column(
            children: [
              for (final card in cards) ...[
                card,
                const SizedBox(height: Spacing.s),
              ],
            ],
          );
        }
        return Row(
          children: [
            for (var i = 0; i < cards.length; i++) ...[
              if (i > 0) const SizedBox(width: Spacing.m),
              Expanded(child: cards[i]),
            ],
          ],
        );
      },
    );
  }
}
