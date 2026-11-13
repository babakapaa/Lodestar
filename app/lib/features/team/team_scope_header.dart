import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import 'team_labels_context.dart';
import 'team_scope.dart';

class TeamScopeHeader extends StatelessWidget {
  const TeamScopeHeader({required this.scope, super.key});

  final TeamScope scope;

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (scope.context.isNotEmpty)
          Text(
            scope.context,
            style: textTheme.bodyMedium!.copyWith(color: context.palette.muted),
          ),
        Text(scope.title, style: textTheme.headlineSmall),
        const SizedBox(height: Spacing.s),
        Text(
          context.teamLabels.rules(scope),
          style: textTheme.bodyMedium!.copyWith(color: context.palette.muted),
        ),
      ],
    );
  }
}
