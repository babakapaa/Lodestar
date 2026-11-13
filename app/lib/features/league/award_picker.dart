import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'award.dart';

abstract final class AwardPicker {
  static Future<Award?> show(BuildContext context, List<Award> awards) =>
      showDialog<Award>(
        context: context,
        builder: (context) => SimpleDialog(
          title: Text(context.l10n.leagueGrant),
          children: [
            if (awards.isEmpty)
              Padding(
                padding: const EdgeInsets.all(24),
                child: Text(context.l10n.leagueAwardsNone),
              ),
            for (final award in awards)
              SimpleDialogOption(
                key: ValueKey('award-${award.id}'),
                onPressed: () => Navigator.of(context).pop(award),
                child: Text(award.title),
              ),
          ],
        ),
      );
}
