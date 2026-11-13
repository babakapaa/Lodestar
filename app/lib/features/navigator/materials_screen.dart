import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/empty_state.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../study/loaded_list_view.dart';
import 'material_tile.dart';
import 'navigator_providers.dart';

class MaterialsScreen extends ConsumerWidget {
  const MaterialsScreen({required this.courseId, super.key});

  final String courseId;

  Future<void> _add(BuildContext context, WidgetRef ref) async {
    final added = await context.push<bool>(AppRoutes.materialNew(courseId));
    if (added ?? false) {
      ref.invalidate(materialsProvider(courseId));
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final value = ref.watch(materialsProvider(courseId));
    return Scaffold(
      appBar: AppBar(title: Text(l10n.materialsTitle)),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _add(context, ref),
        icon: const Icon(Icons.note_add_outlined),
        label: Text(l10n.materialsAdd),
      ),
      body: SafeArea(
        top: false,
        child: LoadedListView(
          value: value,
          bottomInset: 80,
          onRefresh: () => ref.refresh(materialsProvider(courseId).future),
          builder: (list) => [
            Text(
              l10n.materialsHint,
              style: TextStyle(color: context.palette.muted),
            ),
            const SizedBox(height: Spacing.l),
            if (list.isEmpty)
              EmptyState(
                icon: Icons.library_books_outlined,
                title: l10n.materialsEmpty,
                body: l10n.materialsHint,
              )
            else
              for (final material in list)
                MaterialTile(
                  key: ValueKey(material.id),
                  courseId: courseId,
                  material: material,
                ),
          ],
        ),
      ),
    );
  }
}
