import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../../core/sources/gost_reference.dart';
import '../../core/sources/source.dart';

class SourceTile extends StatelessWidget {
  const SourceTile({required this.source, super.key});

  final Source source;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: const Icon(Icons.menu_book_outlined),
      title: Text(GostReference.format(source)),
      subtitle: Text(
        context.l10n.sourcesUnverified,
        style: TextStyle(color: context.palette.muted),
      ),
      onTap: () => context.push(AppRoutes.source(source.id), extra: source),
    );
  }
}
