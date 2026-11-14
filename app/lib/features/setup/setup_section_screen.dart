import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'hardware_card.dart';
import 'legal_question_card.dart';
import 'protocol_card.dart';
import 'section_assignee.dart';
import 'section_link.dart';
import 'section_review_bar.dart';
import 'section_settings.dart';
import 'setup_labels.dart';
import 'setup_overview.dart';
import 'setup_providers.dart';
import 'setup_section.dart';

class SetupSectionScreen extends ConsumerWidget {
  const SetupSectionScreen({required this.sectionKey, super.key});

  final String sectionKey;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final labels = SetupLabels(l10n);
    final manages = ref.watch(setupStatusProvider).value?.manages ?? false;
    final link = SectionLink.of(sectionKey);
    return Scaffold(
      appBar: AppBar(title: Text(labels.section(sectionKey))),
      body: SafeArea(
        top: false,
        child: AsyncListView<SetupOverview>(
          value: ref.watch(setupOverviewProvider),
          onRefresh: () => ref.refresh(setupOverviewProvider.future),
          builder: (overview) {
            final section = overview.section(sectionKey);
            if (section == null) {
              return [Text(labels.hint(sectionKey))];
            }
            return [
              Text(labels.hint(sectionKey)),
              const SizedBox(height: Spacing.s),
              Text(labels.status(section.status)),
              if (section.note.isNotEmpty) ...[
                const SizedBox(height: Spacing.m),
                NoticeBanner(
                  icon: Icons.comment_outlined,
                  message: l10n.setupNote(section.note),
                ),
              ],
              if (manages) SectionAssignee(section: section),
              const SizedBox(height: Spacing.l),
              if (sectionKey == SetupSection.deploy) ...[
                const HardwareCard(),
                const SizedBox(height: Spacing.l),
              ],
              if (sectionKey == SetupSection.legal) ...[
                for (final answer in overview.legal)
                  LegalQuestionCard(
                    key: ValueKey('legal-${answer.question}'),
                    answer: answer,
                  ),
                const SizedBox(height: Spacing.l),
                ProtocolCard(overview: overview),
                const SizedBox(height: Spacing.l),
              ],
              SectionSettings(section: section),
              if (link != null) ...[
                const SizedBox(height: Spacing.m),
                Align(
                  alignment: Alignment.centerLeft,
                  child: LodestarButton(
                    label: l10n.setupOpenScreen,
                    icon: Icons.open_in_new,
                    variant: LodestarButtonVariant.secondary,
                    onPressed: () => context.push(link),
                  ),
                ),
              ],
              const SizedBox(height: Spacing.xl),
              SectionReviewBar(section: section, manages: manages),
            ];
          },
        ),
      ),
    );
  }
}
