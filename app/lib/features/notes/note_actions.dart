import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import 'note.dart';
import 'note_meta.dart';
import 'notes_providers.dart';

abstract final class NoteActions {
  static Future<void> create(
    BuildContext context,
    WidgetRef ref, {
    String? path,
  }) async {
    final notes = ref.read(notesProvider).value ?? const <Note>[];
    final target = path ?? _untitled(context.l10n.notesUntitled, notes);
    final id = await ref
        .read(notesStoreProvider)
        .create(NoteMeta(path: target));
    if (context.mounted) {
      await context.push(AppRoutes.note(id));
    }
  }

  static Future<bool> remote(
    BuildContext context,
    WidgetRef ref,
    Future<void> Function() action,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await action();
      await ref.read(notesSyncStatusProvider.notifier).sync();
      return true;
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
      return false;
    }
  }

  static String _untitled(String base, List<Note> notes) {
    final taken = {for (final note in notes) note.meta.path};
    if (!taken.contains(base)) {
      return base;
    }
    var n = 2;
    while (taken.contains('$base $n')) {
      n++;
    }
    return '$base $n';
  }
}
