import 'package:material_ui/material_ui.dart';

import '../../design/components/choice_group.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'league_display.dart';
import 'league_labels_context.dart';
import 'standing.dart';

class LeagueDisplayDialog extends StatefulWidget {
  const LeagueDisplayDialog({required this.standing, super.key});

  final Standing standing;

  static Future<(LeagueDisplay, String)?> show(
    BuildContext context,
    Standing standing,
  ) => showDialog<(LeagueDisplay, String)>(
    context: context,
    builder: (context) => LeagueDisplayDialog(standing: standing),
  );

  @override
  State<LeagueDisplayDialog> createState() => _LeagueDisplayDialogState();
}

class _LeagueDisplayDialogState extends State<LeagueDisplayDialog> {
  late LeagueDisplay _display = widget.standing.display;
  late final _nickname = TextEditingController(text: widget.standing.nickname);

  @override
  void dispose() {
    _nickname.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(l10n.leagueDisplayTitle),
      content: SizedBox(
        width: 380,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(l10n.leagueDisplayHint),
            const SizedBox(height: Spacing.m),
            ChoiceGroup<LeagueDisplay>(
              values: LeagueDisplay.values,
              selected: _display,
              labelOf: context.leagueLabels.display,
              onSelected: (value) => setState(() => _display = value),
            ),
            if (_display == LeagueDisplay.nickname)
              TextField(
                key: const ValueKey('league-nickname'),
                controller: _nickname,
                maxLength: 40,
                decoration: InputDecoration(labelText: l10n.leagueNickname),
              ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          onPressed: () =>
              Navigator.of(context).pop((_display, _nickname.text.trim())),
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
