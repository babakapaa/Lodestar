import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../l10n/l10n_context.dart';
import 'sandbox_labels_context.dart';
import 'team_pool.dart';

class TeamPoolTile extends StatelessWidget {
  const TeamPoolTile({required this.pool, required this.onGive, super.key});

  final TeamPool pool;
  final VoidCallback onGive;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    return LodestarCard(
      padding: EdgeInsets.zero,
      child: ListTile(
        contentPadding: LodestarCard.tilePadding,
        leading: const Icon(Icons.groups_outlined),
        title: Text(pool.name),
        subtitle: Text(
          [
            labels.hours(pool.available),
            if (pool.reserved > 0)
              l10n.sandboxReserved(labels.hours(pool.reserved)),
          ].join(', '),
        ),
        trailing: TextButton(
          key: ValueKey('give-${pool.id}'),
          onPressed: onGive,
          child: Text(l10n.sandboxGive),
        ),
      ),
    );
  }
}
