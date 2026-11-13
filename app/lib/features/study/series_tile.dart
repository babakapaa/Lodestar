import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'lesson_series.dart';
import 'study_labels_context.dart';

class SeriesTile extends StatelessWidget {
  const SeriesTile({
    required this.series,
    required this.onEdit,
    required this.onDelete,
    super.key,
  });

  final LessonSeries series;
  final VoidCallback onEdit;
  final VoidCallback onDelete;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    return ListTile(
      key: ValueKey('series-${series.id}'),
      contentPadding: EdgeInsets.zero,
      title: Text(
        '${labels.shortWeekday(series.weekday)} ${series.startTime}, ${labels.kind(series.kind)}',
      ),
      subtitle: Text(
        [
          l10n.seriesLength(series.minutes),
          if (series.room.isNotEmpty) series.room,
          l10n.seriesEveryWeeks(series.everyWeeks),
          '${series.startsOn} – ${series.endsOn}',
        ].join(', '),
      ),
      onTap: onEdit,
      trailing: IconButton(
        tooltip: l10n.seriesDelete,
        icon: const Icon(Icons.delete_outline),
        onPressed: onDelete,
      ),
    );
  }
}
