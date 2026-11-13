import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'lti_context_card.dart';
import 'lti_providers.dart';

class LtiContextScreen extends ConsumerWidget {
  const LtiContextScreen({required this.contextId, super.key});

  final String contextId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    return Scaffold(
      appBar: AppBar(title: Text(l10n.lmsCourseTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(ltiContextProvider(contextId)),
          onRefresh: () => ref.refresh(ltiContextProvider(contextId).future),
          builder: (c) => [
            Text(
              c.instructor
                  ? l10n.lmsCourseTeacherHint
                  : l10n.lmsCourseStudentHint,
              style: muted,
            ),
            const SizedBox(height: Spacing.l),
            LtiContextCard(
              context: c,
              onChanged: () => ref.invalidate(ltiContextProvider(contextId)),
            ),
          ],
        ),
      ),
    );
  }
}
