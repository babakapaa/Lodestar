import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../core/navigation/layout_size.dart';
import '../../core/navigation/tab_page.dart';
import '../../design/components/empty_state.dart';
import '../../l10n/l10n_context.dart';
import '../study/loaded_list_view.dart';
import 'dialog_tile.dart';
import 'navigator_providers.dart';
import 'navigator_workspace.dart';
import 'support_sheet.dart';

class NavigatorScreen extends ConsumerWidget {
  const NavigatorScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final value = ref.watch(dialogsProvider);
    final wide =
        LayoutSize.of(MediaQuery.sizeOf(context).width) == LayoutSize.expanded;
    return TabPage(
      title: l10n.tabNavigator,
      actions: [
        IconButton(
          tooltip: l10n.supportTitle,
          icon: const Icon(Icons.volunteer_activism_outlined),
          onPressed: () => SupportSheet.show(context),
        ),
        IconButton(
          tooltip: l10n.journalTitle,
          icon: const Icon(Icons.history),
          onPressed: () => context.push(AppRoutes.journal),
        ),
      ],
      floatingActionButton: wide
          ? null
          : FloatingActionButton.extended(
              onPressed: () => context.push(AppRoutes.newDialog()),
              icon: const Icon(Icons.add_comment_outlined),
              label: Text(l10n.navigatorNewQuestion),
            ),
      child: wide
          ? const NavigatorWorkspace()
          : LoadedListView(
              value: value,
              bottomInset: 80,
              onRefresh: () => ref.refresh(dialogsProvider.future),
              builder: (dialogs) => [
                if (dialogs.isEmpty)
                  EmptyState(
                    icon: Icons.explore_outlined,
                    title: l10n.emptyNavigatorTitle,
                    body: l10n.emptyNavigatorBody,
                  )
                else
                  for (final dialog in dialogs)
                    DialogTile(key: ValueKey(dialog.id), dialog: dialog),
              ],
            ),
    );
  }
}
