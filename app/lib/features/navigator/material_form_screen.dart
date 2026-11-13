import 'dart:convert';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/files/files_providers.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'material_input.dart';
import 'navigator_providers.dart';

class MaterialFormScreen extends ConsumerStatefulWidget {
  const MaterialFormScreen({required this.courseId, super.key});

  final String courseId;

  @override
  ConsumerState<MaterialFormScreen> createState() => _MaterialFormState();
}

class _MaterialFormState extends ConsumerState<MaterialFormScreen> {
  final _title = TextEditingController();
  final _reference = TextEditingController();
  final _url = TextEditingController();
  final _text = TextEditingController();
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    _title.dispose();
    _reference.dispose();
    _url.dispose();
    _text.dispose();
    super.dispose();
  }

  Future<void> _import() async {
    final files = await ref.read(fileOpenerProvider).open(const ['txt', 'md']);
    if (files.isEmpty || !mounted) {
      return;
    }
    final file = files.first;
    setState(() {
      _text.text = utf8.decode(file.bytes, allowMalformed: true);
      if (_title.text.trim().isEmpty) {
        _title.text = file.name.replaceFirst(RegExp(r'\.(txt|md)$'), '');
      }
    });
  }

  Future<void> _save() async {
    final l10n = context.l10n;
    if (_title.text.trim().isEmpty || _text.text.trim().isEmpty) {
      setState(() => _error = l10n.errorInvalidMaterial);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref
          .read(navigatorApiProvider)
          .addMaterial(
            widget.courseId,
            MaterialInput(
              title: _title.text.trim(),
              reference: _reference.text.trim(),
              url: _url.text.trim(),
              text: _text.text,
            ),
          );
      if (mounted) {
        context.pop(true);
      }
    } on Object catch (error) {
      if (mounted) {
        setState(() {
          _busy = false;
          _error = ApiErrorText.of(l10n, error);
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.materialsAdd)),
      body: SafeArea(
        top: false,
        child: ListView(
          padding: const EdgeInsets.all(Spacing.l),
          children: [
            TextField(
              key: const ValueKey('material-title'),
              controller: _title,
              maxLength: 300,
              decoration: InputDecoration(labelText: l10n.materialTitle),
            ),
            const SizedBox(height: Spacing.m),
            TextField(
              key: const ValueKey('material-reference'),
              controller: _reference,
              decoration: InputDecoration(labelText: l10n.materialReference),
            ),
            const SizedBox(height: Spacing.m),
            TextField(
              key: const ValueKey('material-url'),
              controller: _url,
              keyboardType: TextInputType.url,
              decoration: InputDecoration(labelText: l10n.materialUrl),
            ),
            const SizedBox(height: Spacing.m),
            Align(
              alignment: AlignmentDirectional.centerStart,
              child: LodestarButton(
                label: l10n.materialImport,
                icon: Icons.upload_file_outlined,
                variant: LodestarButtonVariant.secondary,
                onPressed: _busy ? null : _import,
              ),
            ),
            const SizedBox(height: Spacing.m),
            TextField(
              key: const ValueKey('material-text'),
              controller: _text,
              minLines: 8,
              maxLines: 20,
              decoration: InputDecoration(
                labelText: l10n.materialText,
                alignLabelWithHint: true,
              ),
            ),
            if (error != null) FormError(message: error),
            const SizedBox(height: Spacing.xl),
            LodestarButton(
              label: l10n.actionSave,
              onPressed: _busy ? null : _save,
            ),
          ],
        ),
      ),
    );
  }
}
