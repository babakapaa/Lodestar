import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';

class RecordModeDialog extends StatelessWidget {
  const RecordModeDialog({super.key});

  static Future<bool?> show(BuildContext context) => showDialog<bool>(
    context: context,
    builder: (context) => const RecordModeDialog(),
  );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return SimpleDialog(
      title: Text(l10n.callRecordTitle),
      children: [
        SimpleDialogOption(
          key: const ValueKey('record-video'),
          onPressed: () => Navigator.of(context).pop(false),
          child: Text(l10n.callRecordVideo),
        ),
        SimpleDialogOption(
          key: const ValueKey('record-audio'),
          onPressed: () => Navigator.of(context).pop(true),
          child: Text(l10n.callRecordAudio),
        ),
      ],
    );
  }
}
