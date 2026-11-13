import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'brief.dart';

class BriefBanner extends StatelessWidget {
  const BriefBanner({required this.brief, super.key});

  final Brief brief;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final textTheme = Theme.of(context).textTheme;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          brief.isCourse
              ? l10n.navigatorAboutCourse(brief.title)
              : l10n.navigatorAboutAssignment(brief.title),
          style: textTheme.titleSmall,
        ),
        const SizedBox(height: Spacing.xs),
        Text(
          brief.project ? l10n.navigatorModeProject : l10n.navigatorModeStudy,
          style: textTheme.bodySmall!.copyWith(color: context.palette.muted),
        ),
      ],
    );
  }
}
