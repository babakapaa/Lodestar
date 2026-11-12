import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'notes_providers.dart';
import 'source.dart';

class SourceFindDialog extends ConsumerStatefulWidget {
  const SourceFindDialog({super.key});

  static Future<Source?> show(BuildContext context) => showDialog<Source>(
    context: context,
    builder: (context) => const SourceFindDialog(),
  );

  @override
  ConsumerState<SourceFindDialog> createState() => _SourceFindDialogState();
}

class _SourceFindDialogState extends ConsumerState<SourceFindDialog> {
  final _controller = TextEditingController();
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _find() async {
    final text = _controller.text.trim();
    if (text.isEmpty) {
      return;
    }
    final l10n = context.l10n;
    final navigator = Navigator.of(context);
    final isLink = text.startsWith('http') && !text.contains('doi.org/');
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      final api = ref.read(sourcesApiProvider);
      final found = isLink
          ? await api.lookup(url: text)
          : await api.lookup(doi: text);
      navigator.pop(found);
    } on ApiException catch (error) {
      setState(() {
        _busy = false;
        _error = ApiErrorText.of(l10n, error);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    return AlertDialog(
      title: Text(l10n.sourcesFind),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          TextField(
            controller: _controller,
            autofocus: true,
            enabled: !_busy,
            decoration: InputDecoration(labelText: l10n.sourcesFindLabel),
            onSubmitted: (_) => _find(),
          ),
          if (error != null) ...[
            const SizedBox(height: Spacing.m),
            FormError(message: error),
          ],
        ],
      ),
      actions: [
        TextButton(
          onPressed: _busy ? null : () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        TextButton(
          onPressed: _busy ? null : _find,
          child: Text(l10n.sourcesFindAction),
        ),
      ],
    );
  }
}
