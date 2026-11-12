import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/brand_mark.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'app_routes.dart';

class EntryPage extends StatelessWidget {
  const EntryPage({
    required this.title,
    required this.subtitle,
    required this.child,
    super.key,
  });

  final String title;
  final String subtitle;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(
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
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(Spacing.l),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 420),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  const Center(child: BrandMark(showName: true)),
                  const SizedBox(height: Spacing.xl),
                  Text(title, style: textTheme.headlineSmall),
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
        ),
      ),
    );
  }
}
