import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';

class FormError extends StatelessWidget {
  const FormError({required this.message, super.key});

  final String message;

  @override
  Widget build(BuildContext context) {
    final palette = context.tokens.palette;
    return Padding(
      padding: const EdgeInsets.only(top: Spacing.m),
      child: Semantics(
        liveRegion: true,
        child: Text(
          message,
          style: Theme.of(context).textTheme.bodyMedium!
              .copyWith(color: palette.warn),
        ),
      ),
    );
  }
}
