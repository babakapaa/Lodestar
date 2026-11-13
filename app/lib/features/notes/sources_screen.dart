import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/empty_state.dart';
import '../../l10n/l10n_context.dart';
import '../study/loaded_list_view.dart';
import '../../core/sources/gost_reference.dart';
import 'notes_providers.dart';
import '../../core/sources/source.dart';
import 'source_find_dialog.dart';
import 'source_tile.dart';

class SourcesScreen extends ConsumerWidget {
  const SourcesScreen({super.key});

  Future<void> _find(BuildContext context) async {
    final found = await SourceFindDialog.show(context);
    if (found != null && context.mounted) {
      await context.push(AppRoutes.sourceNew, extra: found);
    }
  }

  Future<void> _copy(BuildContext context, List<Source> sources) async {
    final messenger = ScaffoldMessenger.of(context);
    final copied = context.l10n.sourcesCopied;
    await Clipboard.setData(ClipboardData(text: GostReference.list(sources)));
    messenger.showSnackBar(SnackBar(content: Text(copied)));
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final value = ref.watch(sourcesProvider);
    final sources = value.value?.value ?? const <Source>[];
    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.sourcesTitle),
        actions: [
          if (sources.isNotEmpty)
            IconButton(
              tooltip: l10n.sourcesCopyList,
              icon: const Icon(Icons.copy_all_outlined),
              onPressed: () => _copy(context, sources),
            ),
          IconButton(
            tooltip: l10n.sourcesManual,
            icon: const Icon(Icons.edit_note),
            onPressed: () => context.push(AppRoutes.sourceNew),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _find(context),
        icon: const Icon(Icons.add_link),
        label: Text(l10n.sourcesAdd),
      ),
      body: LoadedListView(
        value: value,
        onRefresh: () => ref.refresh(sourcesProvider.future),
        bottomInset: 80,
        builder: (list) => [
          if (list.isEmpty)
            EmptyState(
              icon: Icons.menu_book_outlined,
              title: l10n.sourcesEmptyTitle,
              body: l10n.sourcesEmptyBody,
            )
          else
            for (final source in list)
              SourceTile(key: ValueKey(source.id), source: source),
        ],
      ),
    );
  }
}
