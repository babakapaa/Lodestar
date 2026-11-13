import 'package:material_ui/material_ui.dart';

import 'chat_kind.dart';

abstract final class ChatKindIcon {
  static IconData of(ChatKind kind) => switch (kind) {
    ChatKind.direct => Icons.person_outline,
    ChatKind.group => Icons.groups_outlined,
    ChatKind.channel => Icons.campaign_outlined,
  };
}
