import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../attendance/attendance_providers.dart';
import '../league/league_providers.dart';
import 'attendance_trend_section.dart';
import 'course_topics.dart';
import 'divisions_section.dart';
import 'insights_providers.dart';
import 'topics_section.dart';

class GroupInsightsScreen extends ConsumerWidget {
  const GroupInsightsScreen({required this.courseId, super.key});

  final String courseId;

  Future<CourseTopics> _refresh(WidgetRef ref) {
    ref
      ..invalidate(leagueProvider(('study.course', courseId)))
      ..invalidate(courseJournalProvider(courseId));
    return ref.refresh(courseTopicsProvider(courseId).future);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.insightsTitle)),
      body: SafeArea(
        top: false,
        child: RefreshIndicator(
          onRefresh: () => _refresh(ref),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              Text(l10n.insightsPrivacy),
              const SizedBox(height: Spacing.l),
              TopicsSection(courseId: courseId),
              DivisionsSection(courseId: courseId),
              AttendanceTrendSection(courseId: courseId),
            ],
          ),
        ),
      ),
    );
  }
}
