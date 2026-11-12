import 'package:material_ui/material_ui.dart';

import '../../l10n/generated/app_localizations.dart';

enum AppTab {
  today('/today', Icons.wb_sunny_outlined, Icons.wb_sunny),
  study('/study', Icons.school_outlined, Icons.school),
  navigator('/navigator', Icons.explore_outlined, Icons.explore),
  notes('/notes', Icons.sticky_note_2_outlined, Icons.sticky_note_2),
  chats('/chats', Icons.chat_bubble_outline, Icons.chat_bubble);

  const AppTab(this.path, this.icon, this.selectedIcon);

  final String path;
  final IconData icon;
  final IconData selectedIcon;

  String label(AppLocalizations l10n) => switch (this) {
    AppTab.today => l10n.tabToday,
    AppTab.study => l10n.tabStudy,
    AppTab.navigator => l10n.tabNavigator,
    AppTab.notes => l10n.tabNotes,
    AppTab.chats => l10n.tabChats,
  };
}
