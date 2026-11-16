import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class NavigatorStarters extends StatelessWidget {
  const NavigatorStarters({required this.onPick, super.key});

  final ValueChanged<String> onPick;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final items = [
      (Icons.lightbulb_outline, l10n.navigatorStarterExplain),
      (Icons.edit_outlined, l10n.navigatorStarterEssay),
      (Icons.menu_book_outlined, l10n.navigatorStarterSources),
      (Icons.fact_check_outlined, l10n.navigatorStarterCheck),
    ];
    return LayoutBuilder(
      builder: (context, constraints) {
        final columns = constraints.maxWidth >= 520 ? 2 : 1;
        final cell =
            (constraints.maxWidth - Spacing.m * (columns - 1)) / columns;
        return Wrap(
          spacing: Spacing.m,
          runSpacing: Spacing.m,
          children: [
            for (final (icon, text) in items)
              SizedBox(
                width: cell,
                child: _Starter(
                  icon: icon,
                  text: text,
                  onTap: () => onPick(text),
                ),
              ),
          ],
        );
      },
    );
  }
}

class _Starter extends StatelessWidget {
  const _Starter({required this.icon, required this.text, required this.onTap});

  final IconData icon;
  final String text;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    return Material(
      color: palette.surface,
      shape: RoundedRectangleBorder(
        borderRadius: tokens.borderRadius,
        side: BorderSide(color: palette.border),
      ),
      child: InkWell(
        borderRadius: tokens.borderRadius,
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(Spacing.l),
          child: Row(
            children: [
              Icon(icon, size: 20, color: palette.accent),
              const SizedBox(width: Spacing.m),
              Expanded(
                child: Text(
                  text,
                  style: Theme.of(context).textTheme.bodyMedium,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
