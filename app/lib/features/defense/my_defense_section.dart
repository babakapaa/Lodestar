import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/components/notice_banner.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/study_labels_context.dart';
import 'defense_labels.dart';
import 'defense_providers.dart';
import 'defense_status.dart';
import 'my_defense.dart';
import 'protocol_view.dart';
import 'question_types_hint.dart';
import 'slot_text.dart';

class MyDefenseSection extends ConsumerWidget {
  const MyDefenseSection({required this.assignmentId, super.key});

  final String assignmentId;

  Future<void> _cancel(BuildContext context, WidgetRef ref) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await ref.read(defenseApiProvider).cancel(assignmentId);
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(myDefenseProvider(assignmentId));
  }

  List<Widget> _body(BuildContext context, WidgetRef ref, MyDefense my) {
    final l10n = context.l10n;
    final locale = Localizations.localeOf(context).toLanguageTag();
    final slot = my.slot;
    if (my.status.finished) {
      return [
        Text(DefenseLabels(l10n).status(my.status)),
        if (my.protocol case final protocol?
            when my.status == DefenseStatus.held) ...[
          const SizedBox(height: Spacing.m),
          ProtocolView(defense: protocol),
        ],
      ];
    }
    if (!my.required) {
      return [Text(l10n.defenseNotChosen)];
    }
    return [
      if (my.recorded) ...[
        NoticeBanner(
          icon: Icons.videocam_outlined,
          message: l10n.defenseRecorded,
        ),
        const SizedBox(height: Spacing.m),
      ],
      if (slot != null)
        Text(l10n.defenseBooked(SlotText.of(slot, locale)))
      else
        Text(l10n.defenseNotBooked),
      const SizedBox(height: Spacing.m),
      Wrap(
        spacing: Spacing.s,
        runSpacing: Spacing.s,
        children: [
          FilledButton.tonal(
            key: const ValueKey('defense-book'),
            onPressed: () =>
                context.push(AppRoutes.defenseBooking(assignmentId)),
            child: Text(slot == null ? l10n.defenseBook : l10n.defenseRebook),
          ),
          if (slot != null)
            TextButton(
              onPressed: () => _cancel(context, ref),
              child: Text(l10n.defenseCancel),
            ),
        ],
      ),
      const QuestionTypesHint(),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final my = ref.watch(myDefenseProvider(assignmentId));
    final value = my.value;
    return Section(
      title: value == null
          ? l10n.defenseTitle
          : '${l10n.defenseTitle}: ${context.studyLabels.defenseFormat(value.format).toLowerCase()}',
      child: LodestarCard(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: my.hasError
              ? [FormError(message: ApiErrorText.of(l10n, my.error!))]
              : value == null
              ? const [Center(child: CircularProgressIndicator())]
              : _body(context, ref, value),
        ),
      ),
    );
  }
}
