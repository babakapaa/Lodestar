import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class EntryHero extends StatelessWidget {
  const EntryHero({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = context.tokens.palette;
    final textTheme = Theme.of(context).textTheme;
    final ink = palette.onAccent;
    final deep = Color.lerp(palette.accent, const Color(0xFF000000), 0.45)!;
    final features = [
      (Icons.school_outlined, l10n.entryFeatureStudy),
      (Icons.hub_outlined, l10n.entryFeatureNotes),
      (Icons.explore_outlined, l10n.entryFeatureNavigator),
    ];
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [palette.accent, deep],
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.all(Spacing.xxl * 1.5),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: ink.withValues(alpha: 0.16),
                    borderRadius: context.tokens.borderRadiusSmall,
                  ),
                  child: Icon(Icons.explore, color: ink, size: 22),
                ),
                const SizedBox(width: Spacing.m),
                Text(
                  l10n.appTitle,
                  style: textTheme.titleLarge!.copyWith(color: ink),
                ),
              ],
            ),
            const Spacer(),
            ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 440),
              child: Text(
                l10n.entryTagline,
                style: textTheme.displaySmall!.copyWith(
                  color: ink,
                  height: 1.15,
                ),
              ),
            ),
            const SizedBox(height: Spacing.xxl),
            for (final (icon, text) in features)
              Padding(
                padding: const EdgeInsets.only(bottom: Spacing.m),
                child: Row(
                  children: [
                    Icon(icon, color: ink.withValues(alpha: 0.85), size: 20),
                    const SizedBox(width: Spacing.m),
                    Flexible(
                      child: Text(
                        text,
                        style: textTheme.bodyLarge!.copyWith(
                          color: ink.withValues(alpha: 0.9),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            const Spacer(),
          ],
        ),
      ),
    );
  }
}
