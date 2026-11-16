import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';

class TodayActionCard extends StatefulWidget {
  const TodayActionCard({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
    super.key,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  @override
  State<TodayActionCard> createState() => _TodayActionCardState();
}

class _TodayActionCardState extends State<TodayActionCard> {
  bool _hover = false;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    final textTheme = Theme.of(context).textTheme;
    return MouseRegion(
      onEnter: (_) => setState(() => _hover = true),
      onExit: (_) => setState(() => _hover = false),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 150),
        decoration: BoxDecoration(
          color: palette.surface,
          borderRadius: tokens.borderRadius,
          border: Border.all(color: _hover ? palette.accent : palette.border),
          boxShadow: palette.raised,
        ),
        child: Material(
          type: MaterialType.transparency,
          child: InkWell(
            borderRadius: tokens.borderRadius,
            onTap: widget.onTap,
            child: Padding(
              padding: const EdgeInsets.all(Spacing.l),
              child: Row(
                children: [
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: palette.accentSoft,
                      borderRadius: tokens.borderRadiusSmall,
                    ),
                    child: Icon(widget.icon, color: palette.accent, size: 22),
                  ),
                  const SizedBox(width: Spacing.m),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          widget.title,
                          style: textTheme.titleSmall,
                          overflow: TextOverflow.ellipsis,
                        ),
                        const SizedBox(height: 2),
                        Text(
                          widget.subtitle,
                          style: textTheme.bodySmall,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                  Icon(
                    Icons.arrow_forward,
                    size: 18,
                    color: _hover ? palette.accent : palette.muted,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
