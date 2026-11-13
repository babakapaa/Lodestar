import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'sandbox_template.dart';

class TemplateDialog extends StatefulWidget {
  const TemplateDialog({this.value, super.key});

  final SandboxTemplate? value;

  static Future<SandboxTemplate?> show(
    BuildContext context, [
    SandboxTemplate? value,
  ]) => showDialog<SandboxTemplate>(
    context: context,
    builder: (context) => TemplateDialog(value: value),
  );

  @override
  State<TemplateDialog> createState() => _TemplateDialogState();
}

class _TemplateDialogState extends State<TemplateDialog> {
  late final _name = TextEditingController(text: widget.value?.name ?? '');
  late final _about = TextEditingController(text: widget.value?.about ?? '');
  late final _image = TextEditingController(text: widget.value?.image ?? '');
  late bool _starter = widget.value?.starter ?? false;
  late bool _archived = widget.value?.archived ?? false;

  @override
  void dispose() {
    _name.dispose();
    _about.dispose();
    _image.dispose();
    super.dispose();
  }

  void _done() {
    if (_name.text.trim().isEmpty || _image.text.trim().isEmpty) {
      return;
    }
    Navigator.of(context).pop(
      SandboxTemplate(
        id: widget.value?.id ?? '',
        name: _name.text.trim(),
        about: _about.text.trim(),
        image: _image.text.trim(),
        archived: _archived,
        starter: _starter,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(
        widget.value == null
            ? l10n.sandboxTemplateNew
            : l10n.sandboxTemplateEdit,
      ),
      content: SizedBox(
        width: 420,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                key: const ValueKey('template-name'),
                controller: _name,
                decoration: InputDecoration(
                  labelText: l10n.sandboxTemplateName,
                ),
              ),
              TextField(
                key: const ValueKey('template-image'),
                controller: _image,
                decoration: InputDecoration(
                  labelText: l10n.sandboxImage,
                  helperText: l10n.sandboxImageHint,
                ),
              ),
              TextField(
                key: const ValueKey('template-about'),
                controller: _about,
                maxLines: 3,
                decoration: InputDecoration(labelText: l10n.sandboxAbout),
              ),
              SwitchListTile(
                key: const ValueKey('template-starter'),
                contentPadding: EdgeInsets.zero,
                title: Text(l10n.sandboxStarter),
                subtitle: Text(l10n.sandboxStarterHint),
                value: _starter,
                onChanged: (v) => setState(() => _starter = v),
              ),
              SwitchListTile(
                contentPadding: EdgeInsets.zero,
                title: Text(l10n.sandboxArchived),
                value: _archived,
                onChanged: (v) => setState(() => _archived = v),
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
          key: const ValueKey('template-save'),
          onPressed: _done,
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
