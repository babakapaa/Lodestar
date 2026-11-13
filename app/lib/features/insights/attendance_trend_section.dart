import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../attendance/attendance_providers.dart';
import '../attendance/journal_lesson.dart';
import '../study/study_labels_context.dart';
import 'share_bar.dart';

class AttendanceTrendSection extends ConsumerWidget {
  const AttendanceTrendSection({required this.courseId, super.key});

  final String courseId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (!ref.watch(attendanceEnabledProvider)) {
      return const SizedBox.shrink();
    }
    final journal = ref.watch(courseJournalProvider(courseId)).value;
    final held = <JournalLesson>[
      for (final l in journal?.lessons ?? const <JournalLesson>[])
        if (l.sessionId != null && !l.open) l,
    ];
    if (held.isEmpty) {
      return const SizedBox.shrink();
    }
    final l10n = context.l10n;
    final days = context.studyLabels;
    return Padding(
      padding: const EdgeInsets.only(top: Spacing.xl),
      child: Section(
        title: l10n.insightsAttendance,
        child: Column(
          children: [
            for (final l in held.reversed.take(10))
              Builder(
                builder: (context) {
                  final counted = l.counts.here + l.counts.absent;
                  final share = counted == 0 ? 0.0 : l.counts.here / counted;
                  return ShareBar(
                    key: ValueKey('trend-${l.seriesId}-${l.date}'),
                    label: days.shortDay(l.startsAt),
                    value: l10n.attendanceRate((share * 100).round()),
                    share: share,
                  );
                },
              ),
          ],
        ),
      ),
    );
  }
}
