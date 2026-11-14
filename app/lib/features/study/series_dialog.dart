import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'lesson_series.dart';
import 'study_labels_context.dart';

class SeriesDialog extends StatefulWidget {
  const SeriesDialog({this.series, super.key});

  static const kinds = [
    'lecture',
    'practice',
    'lab',
    'seminar',
    'exam',
    'consultation',
  ];

  final LessonSeries? series;

  static Future<LessonSeries?> show(
    BuildContext context, {
    LessonSeries? series,
  }) => showDialog<LessonSeries>(
    context: context,
    builder: (_) => SeriesDialog(series: series),
  );

  @override
  State<SeriesDialog> createState() => _SeriesDialogState();
}

class _SeriesDialogState extends State<SeriesDialog> {
  static final _day = DateFormat('yyyy-MM-dd');

  late String _kind = widget.series?.kind ?? SeriesDialog.kinds.first;
  late int _weekday = widget.series?.weekday ?? DateTime.now().weekday;
  late int _every = widget.series?.everyWeeks ?? 1;
  late DateTime _from = _parse(widget.series?.startsOn) ?? DateTime.now();
  late DateTime _to =
      _parse(widget.series?.endsOn) ??
      DateTime.now().add(const Duration(days: 120));
  late final _time = TextEditingController(
    text: widget.series?.startTime ?? '09:00',
  );
  late final _minutes = TextEditingController(
    text: '${widget.series?.minutes ?? 90}',
  );
  late final _room = TextEditingController(text: widget.series?.room ?? '');
  bool _invalid = false;

  static DateTime? _parse(String? day) =>
      day == null ? null : DateTime.tryParse(day);

  @override
  void dispose() {
    _time.dispose();
    _minutes.dispose();
    _room.dispose();
    super.dispose();
  }

  Future<void> _pick({required bool start}) async {
    final picked = await showDatePicker(
      context: context,
      initialDate: start ? _from : _to,
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
    );
    if (picked != null) {
      setState(() => start ? _from = picked : _to = picked);
    }
  }

  void _save() {
    final minutes = int.tryParse(_minutes.text.trim());
    final time = _time.text.trim();
    if (minutes == null ||
        minutes < 1 ||
        !RegExp(r'^([01]\d|2[0-3]):[0-5]\d$').hasMatch(time) ||
        _to.isBefore(_from)) {
      setState(() => _invalid = true);
      return;
    }
    Navigator.of(context).pop(
      LessonSeries(
        id: widget.series?.id ?? '',
        kind: _kind,
        room: _room.text.trim(),
        startsOn: _day.format(_from),
        endsOn: _day.format(_to),
        weekday: _weekday,
        startTime: time,
        minutes: minutes,
        everyWeeks: _every,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    return AlertDialog(
      title: Text(widget.series == null ? l10n.seriesAdd : l10n.seriesEdit),
      content: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            DropdownButtonFormField<String>(
              isExpanded: true,
              initialValue: _kind,
              decoration: InputDecoration(labelText: l10n.seriesKind),
              items: [
                for (final kind in SeriesDialog.kinds)
                  DropdownMenuItem(value: kind, child: Text(labels.kind(kind))),
              ],
              onChanged: (value) => setState(() => _kind = value ?? _kind),
            ),
            const SizedBox(height: Spacing.m),
            DropdownButtonFormField<int>(
              isExpanded: true,
              initialValue: _weekday,
              decoration: InputDecoration(labelText: l10n.seriesWeekday),
              items: [
                for (var day = 1; day <= 7; day++)
                  DropdownMenuItem(
                    value: day,
                    child: Text(labels.weekday(day)),
                  ),
              ],
              onChanged: (value) =>
                  setState(() => _weekday = value ?? _weekday),
            ),
            const SizedBox(height: Spacing.m),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    key: const ValueKey('series-time'),
                    controller: _time,
                    decoration: InputDecoration(labelText: l10n.seriesTime),
                  ),
                ),
                const SizedBox(width: Spacing.m),
                Expanded(
                  child: TextField(
                    key: const ValueKey('series-minutes'),
                    controller: _minutes,
                    keyboardType: TextInputType.number,
                    decoration: InputDecoration(labelText: l10n.seriesMinutes),
                  ),
                ),
              ],
            ),
            const SizedBox(height: Spacing.m),
            TextField(
              key: const ValueKey('series-room'),
              controller: _room,
              decoration: InputDecoration(labelText: l10n.seriesRoom),
            ),
            const SizedBox(height: Spacing.m),
            DropdownButtonFormField<int>(
              isExpanded: true,
              initialValue: _every,
              decoration: InputDecoration(labelText: l10n.seriesEvery),
              items: [
                for (var weeks = 1; weeks <= 4; weeks++)
                  DropdownMenuItem(
                    value: weeks,
                    child: Text(l10n.seriesEveryWeeks(weeks)),
                  ),
              ],
              onChanged: (value) => setState(() => _every = value ?? _every),
            ),
            const SizedBox(height: Spacing.s),
            ListTile(
              contentPadding: EdgeInsets.zero,
              title: Text(l10n.seriesFrom),
              trailing: Text(labels.date(_from)),
              onTap: () => _pick(start: true),
            ),
            ListTile(
              contentPadding: EdgeInsets.zero,
              title: Text(l10n.seriesTo),
              trailing: Text(labels.date(_to)),
              onTap: () => _pick(start: false),
            ),
            if (_invalid)
              Text(
                l10n.errorInvalidSeries,
                style: TextStyle(color: Theme.of(context).colorScheme.error),
              ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          key: const ValueKey('series-save'),
          onPressed: _save,
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
