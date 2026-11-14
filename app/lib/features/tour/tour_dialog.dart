import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'tour_dots.dart';
import 'tour_steps.dart';

class TourDialog extends StatefulWidget {
  const TourDialog({super.key});

  static Future<void> show(BuildContext context) => showDialog<void>(
    context: context,
    barrierDismissible: false,
    builder: (context) => const TourDialog(),
  );

  @override
  State<TourDialog> createState() => _TourDialogState();
}

class _TourDialogState extends State<TourDialog> {
  int _index = 0;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final tokens = context.tokens;
    final textTheme = Theme.of(context).textTheme;
    final steps = TourSteps.of(l10n);
    final step = steps[_index];
    final last = _index == steps.length - 1;
    return Dialog(
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 440),
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
                  child: Icon(
                    step.icon,
                    color: tokens.palette.accent,
                    size: 32,
                  ),
                ),
              ),
              const SizedBox(height: Spacing.l),
              Text(
                step.title,
                key: const ValueKey('tour-title'),
                style: textTheme.titleLarge,
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: Spacing.s),
              Text(
                step.body,
                style: textTheme.bodyMedium!.copyWith(
                  color: tokens.palette.muted,
                  height: 1.5,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: Spacing.xl),
              TourDots(count: steps.length, current: _index),
              const SizedBox(height: Spacing.xl),
              Row(
                children: [
                  if (!last)
                    TextButton(
                      onPressed: () => Navigator.of(context).pop(),
                      child: Text(l10n.tourSkip),
                    ),
                  const Spacer(),
                  if (_index > 0)
                    TextButton(
                      onPressed: () => setState(() => _index--),
                      child: Text(l10n.tourBack),
                    ),
                  const SizedBox(width: Spacing.s),
                  FilledButton(
                    onPressed: last
                        ? () => Navigator.of(context).pop()
                        : () => setState(() => _index++),
                    child: Text(last ? l10n.tourDone : l10n.tourNext),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
