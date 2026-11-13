import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'chat_message.dart';

abstract final class ChatPreview {
  static String of(BuildContext context, ChatMessage? message) {
    final l10n = context.l10n;
    if (message == null) {
      return l10n.chatsEmpty;
    }
    if (message.deleted) {
      return l10n.chatsDeleted;
    }
    if (message.text.isEmpty && message.files.isNotEmpty) {
      return '${l10n.chatsFile}: ${message.files.first.name}';
    }
    return message.text.replaceAll('\n', ' ');
  }
}
