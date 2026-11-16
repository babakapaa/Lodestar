import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/brand_mark.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'app_routes.dart';
import 'entry_hero.dart';

class EntryPage extends StatelessWidget {
  const EntryPage({
    required this.title,
    required this.subtitle,
    required this.child,
    super.key,
  });

  static const wideFrom = 900.0;

  final String title;
  final String subtitle;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        shape: const Border(),
        actions: [
          IconButton(
            tooltip: context.l10n.settingsTitle,
            onPressed: () => context.push(AppRoutes.settings),
            icon: const Icon(Icons.tune),
          ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: LayoutBuilder(
          builder: (context, constraints) {
            final form = Center(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(Spacing.l),
                child: ConstrainedBox(
                  constraints: const BoxConstraints(maxWidth: 420),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      if (constraints.maxWidth < wideFrom) ...[
                        const Center(child: BrandMark(showName: true)),
                        const SizedBox(height: Spacing.xl),
                      ],
                      Text(title, style: textTheme.headlineMedium),
                      const SizedBox(height: Spacing.s),
                      Text(
                        subtitle,
                        style: textTheme.bodyMedium!.copyWith(
                          color: context.tokens.palette.muted,
                        ),
                      ),
                      const SizedBox(height: Spacing.xl),
                      LodestarCard(child: child),
                    ],
                  ),
                ),
              ),
            );
            if (constraints.maxWidth < wideFrom) {
              return form;
            }
            return Row(
              children: [
                const Expanded(child: EntryHero()),
                Expanded(child: form),
              ],
            );
          },
        ),
      ),
    );
  }
}
