import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class SupportCard extends StatelessWidget {
  const SupportCard({required this.contacts, super.key});

  final String contacts;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    return Semantics(
      liveRegion: true,
      child: DecoratedBox(
        decoration: BoxDecoration(
          color: tokens.palette.warnSoft,
          borderRadius: tokens.borderRadius,
        ),
        child: Padding(
          padding: const EdgeInsets.all(Spacing.m),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Icon(Icons.support_outlined, color: tokens.palette.warn),
                  const SizedBox(width: Spacing.s),
                  Expanded(
                    child: Text(
                      context.l10n.navigatorSupport,
                      style: Theme.of(context).textTheme.labelLarge,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: Spacing.s),
              SelectableText(contacts),
            ],
          ),
        ),
      ),
    );
  }
}
