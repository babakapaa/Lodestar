import 'package:flutter/services.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_button.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class LetterSheet extends StatefulWidget {
  const LetterSheet({required this.draft, super.key});

  final String draft;

  static Future<void> show(BuildContext context, String draft) {
    return showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      builder: (context) => LetterSheet(draft: draft),
    );
  }

  @override
  State<LetterSheet> createState() => _LetterSheetState();
}

class _LetterSheetState extends State<LetterSheet> {
  late final _text = TextEditingController(text: widget.draft);

  @override
  void dispose() {
    _text.dispose();
    super.dispose();
  }

  Future<void> _copy() async {
    final messenger = ScaffoldMessenger.of(context);
    final copied = context.l10n.letterCopied;
    await Clipboard.setData(ClipboardData(text: _text.text));
    messenger.showSnackBar(SnackBar(content: Text(copied)));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return SafeArea(
      child: Padding(
        padding: EdgeInsets.fromLTRB(
          Spacing.l,
          0,
          Spacing.l,
          Spacing.l + MediaQuery.viewInsetsOf(context).bottom,
        ),
        child: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                l10n.letterTitle,
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: Spacing.s),
              Text(
                l10n.letterHint,
                style: TextStyle(color: context.palette.muted),
              ),
              const SizedBox(height: Spacing.m),
              TextField(
                key: const ValueKey('letter-text'),
                controller: _text,
                minLines: 8,
                maxLines: 16,
              ),
              const SizedBox(height: Spacing.l),
              LodestarButton(
                label: l10n.letterCopy,
                icon: Icons.copy_outlined,
                onPressed: _copy,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
