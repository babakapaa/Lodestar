import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_text_dialog.dart';
import 'lesson.dart';
import 'lesson_action.dart';
import 'lesson_actions_sheet.dart';
import 'lesson_change.dart';
import 'lesson_series.dart';
import 'lesson_tile.dart';
import 'schedule_providers.dart';
import 'series_dialog.dart';
import 'series_tile.dart';
import '../../design/components/tile_group.dart';

class CourseScheduleScreen extends ConsumerWidget {
  const CourseScheduleScreen({required this.courseId, super.key});

  final String courseId;

  Future<void> _run(
    BuildContext context,
    WidgetRef ref,
    Future<void> Function() action,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await action();
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref
      ..invalidate(courseSeriesProvider(courseId))
      ..invalidate(upcomingLessonsProvider(courseId));
  }

  Future<void> _edit(
    BuildContext context,
    WidgetRef ref, [
    LessonSeries? series,
  ]) async {
    final saved = await SeriesDialog.show(context, series: series);
    if (saved == null || !context.mounted) {
      return;
    }
    final api = ref.read(scheduleApiProvider);
    await _run(
      context,
      ref,
      () => series == null ? api.create(courseId, saved) : api.update(saved),
    );
  }

  Future<void> _lesson(
    BuildContext context,
    WidgetRef ref,
    Lesson lesson,
  ) async {
    final l10n = context.l10n;
    final action = await LessonActionsSheet.show(context, lesson);
    if (action == null || !context.mounted) {
      return;
    }
    final api = ref.read(scheduleApiProvider);
    switch (action) {
      case LessonAction.cancel:
        final note = await TeamTextDialog.show(
          context,
          title: l10n.lessonCancel,
          label: l10n.lessonNote,
          action: l10n.lessonCancel,
          required: false,
        );
        if (note != null && context.mounted) {
          await _run(
            context,
            ref,
            () => api.change(
              lesson.seriesId,
              lesson.date,
              LessonChange.cancel(note: note),
            ),
          );
        }
      case LessonAction.room:
        final room = await TeamTextDialog.show(
          context,
          title: l10n.lessonRoom,
          label: l10n.seriesRoom,
          action: l10n.actionSave,
          initial: lesson.room,
        );
        if (room != null && context.mounted) {
          await _run(
            context,
            ref,
            () => api.change(
              lesson.seriesId,
              lesson.date,
              LessonChange.move(room: room),
            ),
          );
        }
      case LessonAction.move:
        final start = lesson.startsAt.toLocal();
        final day = await showDatePicker(
          context: context,
          initialDate: start,
          firstDate: start.subtract(const Duration(days: 365)),
          lastDate: start.add(const Duration(days: 365)),
        );
        if (day == null || !context.mounted) {
          return;
        }
        final time = await showTimePicker(
          context: context,
          initialTime: TimeOfDay.fromDateTime(start),
        );
        if (time == null || !context.mounted) {
          return;
        }
        final at = DateTime(
          day.year,
          day.month,
          day.day,
          time.hour,
          time.minute,
        );
        await _run(
          context,
          ref,
          () => api.change(
            lesson.seriesId,
            lesson.date,
            LessonChange.move(startsAt: at),
          ),
        );
      case LessonAction.revert:
        await _run(
          context,
          ref,
          () => api.revert(lesson.seriesId, lesson.date),
        );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final series = ref.watch(courseSeriesProvider(courseId));
    final lessons = ref.watch(upcomingLessonsProvider(courseId));
    return Scaffold(
      appBar: AppBar(title: Text(l10n.scheduleEditTitle)),
      floatingActionButton: FloatingActionButton.extended(
        key: const ValueKey('series-add'),
        onPressed: () => _edit(context, ref),
        icon: const Icon(Icons.add),
        label: Text(l10n.seriesAdd),
      ),
      body: SafeArea(
        top: false,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(
            Spacing.l,
            Spacing.l,
            Spacing.l,
            96,
          ),
          children: [
            Section(
              title: l10n.seriesTitle,
              child: switch (series) {
                AsyncData(:final value) when value.isEmpty => Text(
                  l10n.seriesEmpty,
                ),
                AsyncData(:final value) => Column(
                  children: [
                    for (final item in value)
                      SeriesTile(
                        series: item,
                        onEdit: () => _edit(context, ref, item),
                        onDelete: () => _run(
                          context,
                          ref,
                          () => ref.read(scheduleApiProvider).delete(item.id),
                        ),
                      ),
                  ],
                ),
                AsyncError(:final error) => Text(ApiErrorText.of(l10n, error)),
                _ => const LinearProgressIndicator(),
              },
            ),
            const SizedBox(height: Spacing.xl),
            Section(
              title: l10n.lessonsUpcoming,
              child: switch (lessons) {
                AsyncData(:final value) when value.isEmpty => Text(
                  l10n.lessonsNone,
                ),
                AsyncData(:final value) => TileGroup(
                  children: [
                    for (final lesson in value)
                      LessonTile(
                        key: ValueKey(
                          'lesson-${lesson.seriesId}-${lesson.date}',
                        ),
                        lesson: lesson,
                        onTap: () => _lesson(context, ref, lesson),
                      ),
                  ],
                ),
                AsyncError(:final error) => Text(ApiErrorText.of(l10n, error)),
                _ => const LinearProgressIndicator(),
              },
            ),
          ],
        ),
      ),
    );
  }
}
