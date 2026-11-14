import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';

class NavigatorTile extends StatelessWidget {
  const NavigatorTile({
    required this.icon,
    required this.title,
    required this.route,
    this.subtitle,
    super.key,
  });

  final IconData icon;
  final String title;
  final String? subtitle;
  final String route;

  @override
  Widget build(BuildContext context) {
    final text = subtitle;
    return LodestarCard(
      padding: EdgeInsets.zero,
      child: ListTile(
        contentPadding: LodestarCard.tilePadding,
        leading: Icon(icon),
        title: Text(title),
        subtitle: text == null ? null : Text(text),
        trailing: const Icon(Icons.chevron_right),
        onTap: () => context.push(route),
      ),
    );
  }
}
