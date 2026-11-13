import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'award.dart';
import 'award_input.dart';

class AwardDialog extends StatefulWidget {
  const AwardDialog({this.award, super.key});

  final Award? award;

  static Future<AwardInput?> show(BuildContext context, [Award? award]) =>
      showDialog<AwardInput>(
        context: context,
        builder: (context) => AwardDialog(award: award),
      );

  @override
  State<AwardDialog> createState() => _AwardDialogState();
}

class _AwardDialogState extends State<AwardDialog> {
  late final _title = TextEditingController(text: widget.award?.title ?? '');
  late final _about = TextEditingController(text: widget.award?.about ?? '');
  late final _gpu = TextEditingController(
    text: '${widget.award?.gpuHours ?? 0}',
  );
  late final _quantity = TextEditingController(
    text: widget.award?.quantity?.toString() ?? '',
  );
  late bool _archived = widget.award?.archived ?? false;

  @override
  void dispose() {
    for (final c in [_title, _about, _gpu, _quantity]) {
      c.dispose();
    }
    super.dispose();
  }

  void _done() {
    final gpu = int.tryParse(_gpu.text.trim());
    final quantityText = _quantity.text.trim();
    final quantity = quantityText.isEmpty ? null : int.tryParse(quantityText);
    if (_title.text.trim().isEmpty ||
        gpu == null ||
        gpu < 0 ||
        (quantityText.isNotEmpty && (quantity == null || quantity < 1))) {
      return;
    }
    Navigator.of(context).pop(
      AwardInput(
        title: _title.text.trim(),
        about: _about.text.trim(),
        gpuHours: gpu,
        quantity: quantity,
        archived: _archived,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(widget.award == null ? l10n.awardNew : l10n.awardEdit),
      content: SizedBox(
        width: 380,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                key: const ValueKey('award-title'),
                controller: _title,
                decoration: InputDecoration(labelText: l10n.awardTitle),
              ),
              TextField(
                controller: _about,
                maxLines: 4,
                decoration: InputDecoration(labelText: l10n.awardAbout),
              ),
              TextField(
                key: const ValueKey('award-gpu'),
                controller: _gpu,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(labelText: l10n.awardGpu),
              ),
              TextField(
                key: const ValueKey('award-quantity'),
                controller: _quantity,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(
                  labelText: l10n.awardQuantity,
                  helperText: l10n.awardQuantityHint,
                ),
              ),
              if (widget.award != null)
                SwitchListTile(
                  contentPadding: EdgeInsets.zero,
                  title: Text(l10n.awardArchived),
                  value: _archived,
                  onChanged: (value) => setState(() => _archived = value),
                ),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          key: const ValueKey('award-save'),
          onPressed: _done,
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
