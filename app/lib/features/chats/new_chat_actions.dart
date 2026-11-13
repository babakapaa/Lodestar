import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import '../people/person_picker_dialog.dart';
import 'chat_details.dart';
import 'chat_kind.dart';
import 'chat_kind_icon.dart';
import 'chat_title_dialog.dart';
import 'chats_providers.dart';
import 'new_chat.dart';

abstract final class NewChatActions {
  static Future<void> choose(BuildContext context, WidgetRef ref) async {
    final l10n = context.l10n;
    final kind = await showModalBottomSheet<ChatKind>(
      context: context,
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            for (final (kind, label) in [
              (ChatKind.direct, l10n.chatsNewDirect),
              (ChatKind.group, l10n.chatsNewGroup),
              (ChatKind.channel, l10n.chatsNewChannel),
            ])
              ListTile(
                leading: Icon(ChatKindIcon.of(kind)),
                title: Text(label),
                onTap: () => Navigator.of(context).pop(kind),
              ),
          ],
        ),
      ),
    );
    if (kind == null || !context.mounted) {
      return;
    }
    await (kind == ChatKind.direct
        ? direct(context, ref)
        : create(context, ref, kind));
  }

  static Future<void> direct(BuildContext context, WidgetRef ref) async {
    final person = await PersonPickerDialog.show(context);
    if (person == null || !context.mounted) {
      return;
    }
    await open(
      context,
      ref,
      () => ref.read(chatsApiProvider).direct(person.id),
    );
  }

  static Future<void> create(
    BuildContext context,
    WidgetRef ref,
    ChatKind kind,
  ) async {
    final l10n = context.l10n;
    final title = await ChatTitleDialog.show(
      context,
      kind == ChatKind.channel ? l10n.chatsNewChannel : l10n.chatsNewGroup,
    );
    if (title == null || !context.mounted) {
      return;
    }
    await open(
      context,
      ref,
      () =>
          ref.read(chatsApiProvider).create(NewChat(kind: kind, title: title)),
    );
  }

  static Future<void> open(
    BuildContext context,
    WidgetRef ref,
    Future<ChatDetails> Function() request,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      final details = await request();
      await ref.read(chatsStoreProvider).saveChat(details.chat);
      if (context.mounted) {
        await context.push(AppRoutes.chat(details.chat.id));
      }
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
  }
}
