import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'notes_providers.dart';
import '../../core/sources/source.dart';
import 'source_edit_screen.dart';

class SourceScreen extends ConsumerWidget {
  const SourceScreen({required this.sourceId, this.source, super.key});

  final String sourceId;
  final Source? source;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final known = source;
    if (known != null) {
      return SourceEditScreen(source: known, existing: true);
    }
    final list = ref.watch(sourcesProvider).value?.value;
    if (list == null) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }
    for (final item in list) {
      if (item.id == sourceId) {
        return SourceEditScreen(source: item, existing: true);
      }
    }
    return Scaffold(
      appBar: AppBar(),
      body: Padding(
        padding: const EdgeInsets.all(Spacing.l),
        child: Text(context.l10n.errorSourceNotFound),
      ),
    );
  }
}
