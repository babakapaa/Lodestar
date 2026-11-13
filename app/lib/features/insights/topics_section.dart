import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/notice_banner.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../navigator/navigator_labels_context.dart';
import '../study/study_labels_context.dart';
import 'course_topics.dart';
import 'insights_providers.dart';
import 'share_bar.dart';

class TopicsSection extends ConsumerWidget {
  const TopicsSection({required this.courseId, super.key});

  final String courseId;

  List<Widget> _body(BuildContext context, CourseTopics topics, int askers) {
    final l10n = context.l10n;
    final suggestion = topics.suggestion;
    final labels = context.navigatorLabels;
    final days = context.studyLabels;
    final most = [
      for (final t in topics.topics) t.students,
      1,
    ].reduce((a, b) => a > b ? a : b);
    return [
      Text(l10n.insightsAskers(askers)),
      const SizedBox(height: Spacing.s),
      if (suggestion != null)
        NoticeBanner(
          key: const ValueKey('insights-suggestion'),
          icon: Icons.lightbulb_outline,
          message: l10n.insightsSuggestion(suggestion.title, suggestion.recent),
        ),
      if (topics.topics.isEmpty)
        Text(l10n.insightsNoTopics(topics.threshold))
      else
        for (final t in topics.topics)
          ShareBar(
            key: ValueKey('topic-${t.materialId}'),
            label: t.title,
            value: l10n.insightsStudents(t.students),
            share: t.students / most,
          ),
      if (topics.intents.isNotEmpty) ...[
        const SizedBox(height: Spacing.l),
        Text(
          l10n.insightsIntents,
          style: Theme.of(context).textTheme.titleSmall,
        ),
        for (final i in topics.intents)
          if (i.intent case final intent?)
            ShareBar(
              label: labels.intent(intent),
              value: l10n.insightsStudents(i.students),
              share: i.students / askers,
            ),
      ],
      if (topics.trend.isNotEmpty) ...[
        const SizedBox(height: Spacing.l),
        Text(l10n.insightsTrend, style: Theme.of(context).textTheme.titleSmall),
        for (final w in topics.trend)
          ShareBar(
            label: l10n.insightsWeekOf(days.shortDay(w.week)),
            value: switch (w.students) {
              final n? => l10n.insightsStudents(n),
              null => l10n.insightsWeekHidden(topics.threshold),
            },
            share: (w.students ?? 0) / askers,
          ),
      ],
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final topics = ref.watch(courseTopicsProvider(courseId));
    return Section(
      title: l10n.insightsTopics,
      child: switch (topics) {
        AsyncData(:final value) => Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: switch (value.students) {
            final askers? => _body(context, value, askers),
            null => [Text(l10n.insightsHidden(value.threshold))],
          },
        ),
        AsyncError() => Text(l10n.insightsUnavailable),
        _ => const Center(child: CircularProgressIndicator()),
      },
    );
  }
}
