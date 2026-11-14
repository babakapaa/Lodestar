import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/empty_state.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'lesson.dart';
import 'lesson_tile.dart';
import 'loaded_list_view.dart';
import 'study_day.dart';
import 'study_labels_context.dart';
import 'study_providers.dart';
import '../../design/components/tile_group.dart';

class ScheduleScreen extends ConsumerStatefulWidget {
  const ScheduleScreen({this.date, super.key});

  final String? date;

  @override
  ConsumerState<ScheduleScreen> createState() => _ScheduleScreenState();
}

class _ScheduleScreenState extends ConsumerState<ScheduleScreen> {
  late DateTime _start = StudyDay.weekStart(StudyDay.parse(widget.date ?? ''));

  (String, String) get _range => (
    StudyDay.format(_start),
    StudyDay.format(_start.add(const Duration(days: 6))),
  );

  void _shift(int weeks) => setState(
    () => _start = DateTime(_start.year, _start.month, _start.day + 7 * weeks),
  );

  List<Widget> _days(BuildContext context, List<Lesson> lessons) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final byDay = <String, List<Lesson>>{};
    for (final lesson in lessons) {
      byDay.putIfAbsent(lesson.date, () => []).add(lesson);
    }
    if (byDay.isEmpty) {
      return [
        EmptyState(
          icon: Icons.event_available_outlined,
          title: l10n.scheduleEmpty,
          body: l10n.scheduleEmptyBody,
        ),
      ];
    }
    return [
      for (final MapEntry(key: day, value: list) in byDay.entries) ...[
        Section(
          title: labels.day(StudyDay.parse(day)),
          child: TileGroup(
            children: [
              for (final lesson in list)
                LessonTile(
                  key: ValueKey('${lesson.seriesId}/${lesson.date}'),
                  lesson: lesson,
                ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.xl),
      ],
    ];
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final range = _range;
    final end = _start.add(const Duration(days: 6));
    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.scheduleTitle),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(48),
          child: Row(
            children: [
              IconButton(
                tooltip: l10n.schedulePrevWeek,
                icon: const Icon(Icons.chevron_left),
                onPressed: () => _shift(-1),
              ),
              Expanded(
                child: Text(
                  '${labels.shortDay(_start)} – ${labels.shortDay(end)}',
                  textAlign: TextAlign.center,
                ),
              ),
              IconButton(
                tooltip: l10n.scheduleNextWeek,
                icon: const Icon(Icons.chevron_right),
                onPressed: () => _shift(1),
              ),
            ],
          ),
        ),
      ),
      body: SafeArea(
        top: false,
        child: LoadedListView(
          value: ref.watch(weekProvider(range)),
          onRefresh: () => ref.refresh(weekProvider(range).future),
          builder: (lessons) => _days(context, lessons),
        ),
      ),
    );
  }
}
