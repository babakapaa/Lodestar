import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';

class EmptyState extends StatelessWidget {
  const EmptyState({
    required this.icon,
    required this.title,
    required this.body,
    super.key,
  });

  final IconData icon;
  final String title;
  final String body;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final textTheme = Theme.of(context).textTheme;
    return Center(
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 360),
        child: Padding(
          padding: const EdgeInsets.all(Spacing.xl),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              DecoratedBox(
                decoration: BoxDecoration(
                  color: tokens.palette.accentSoft,
                  borderRadius: tokens.borderRadius,
                ),
                child: Padding(
                  padding: const EdgeInsets.all(Spacing.l),
                  child: Icon(icon, color: tokens.palette.accent, size: 28),
                ),
              ),
              const SizedBox(height: Spacing.l),
              Text(
                title,
                style: textTheme.titleLarge,
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: Spacing.s),
              Text(
                body,
                style: textTheme.bodyMedium!.copyWith(
                  color: tokens.palette.muted,
                ),
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
