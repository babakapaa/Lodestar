import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import 'navigator_dialog.dart';

class DialogTile extends StatelessWidget {
  const DialogTile({required this.dialog, super.key});

  final NavigatorDialog dialog;

  @override
  Widget build(BuildContext context) {
    final locale = Localizations.localeOf(context).toLanguageTag();
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Icon(
        dialog.scopeKind.isEmpty
            ? Icons.chat_bubble_outline
            : Icons.assignment_outlined,
      ),
      title: Text(
        dialog.title.isEmpty ? context.l10n.navigatorUntitled : dialog.title,
        maxLines: 2,
        overflow: TextOverflow.ellipsis,
      ),
      subtitle: Text(
        DateFormat.yMMMd(locale).add_Hm().format(dialog.changedAt.toLocal()),
      ),
      onTap: () => context.push(AppRoutes.dialog(dialog.id)),
    );
  }
}
