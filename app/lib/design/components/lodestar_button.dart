import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';
import 'lodestar_button_variant.dart';

class LodestarButton extends StatelessWidget {
  const LodestarButton({
    required this.label,
    required this.onPressed,
    this.variant = LodestarButtonVariant.primary,
    this.icon,
    super.key,
  });

  final String label;
  final VoidCallback? onPressed;
  final LodestarButtonVariant variant;
  final IconData? icon;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    final primary = variant == LodestarButtonVariant.primary;
    final style = ButtonStyle(
      minimumSize: const WidgetStatePropertyAll(Size(44, 44)),
      padding: const WidgetStatePropertyAll(
        EdgeInsets.symmetric(horizontal: Spacing.l),
      ),
      shape: WidgetStatePropertyAll(
        RoundedRectangleBorder(borderRadius: tokens.borderRadiusSmall),
      ),
      side: WidgetStatePropertyAll(
        primary ? BorderSide.none : BorderSide(color: palette.border),
      ),
      backgroundColor: WidgetStatePropertyAll(
        primary ? palette.accent : Colors.transparent,
      ),
      foregroundColor: WidgetStatePropertyAll(
        primary ? palette.onAccent : palette.text,
      ),
      textStyle: WidgetStatePropertyAll(Theme.of(context).textTheme.labelLarge),
      elevation: const WidgetStatePropertyAll(0),
    );
    final iconData = icon;
    if (iconData == null) {
      return TextButton(onPressed: onPressed, style: style, child: Text(label));
    }
    return TextButton.icon(
      onPressed: onPressed,
      style: style,
      icon: Icon(iconData, size: 18),
      label: Text(label),
    );
  }
}
