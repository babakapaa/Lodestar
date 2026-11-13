import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/lodestar_card.dart';
import 'team_card.dart';

class TeamCardTile extends StatelessWidget {
  const TeamCardTile({required this.card, this.trailing, super.key});

  final TeamCard card;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    final scope = card.scope;
    return LodestarCard(
      child: ListTile(
        contentPadding: EdgeInsets.zero,
        leading: const Icon(Icons.groups_outlined),
        title: Text(card.team.name),
        subtitle: Text(
          [
            if (scope.context.isNotEmpty) scope.context,
            scope.title,
            card.team.names,
          ].join('\n'),
        ),
        isThreeLine: true,
        trailing: trailing ?? const Icon(Icons.chevron_right),
        onTap: trailing == null
            ? () => context.push(AppRoutes.team(card.team.id))
            : null,
      ),
    );
  }
}
