import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/study_labels_context.dart';
import 'new_slots.dart';

class AddSlotsDialog extends StatefulWidget {
  const AddSlotsDialog({super.key});

  static Future<NewSlots?> show(BuildContext context) => showDialog<NewSlots>(
    context: context,
    builder: (context) => const AddSlotsDialog(),
  );

  @override
  State<AddSlotsDialog> createState() => _AddSlotsDialogState();
}

class _AddSlotsDialogState extends State<AddSlotsDialog> {
  late DateTime _day = DateUtils.dateOnly(
    DateTime.now().add(const Duration(days: 1)),
  );
  TimeOfDay _time = const TimeOfDay(hour: 10, minute: 0);
  final _minutes = TextEditingController(text: '15');
  final _count = TextEditingController(text: '4');
  final _capacity = TextEditingController(text: '1');
  final _place = TextEditingController();
  bool _invalid = false;

  @override
  void dispose() {
    _minutes.dispose();
    _count.dispose();
    _capacity.dispose();
    _place.dispose();
    super.dispose();
  }

  Future<void> _pickDay() async {
    final now = DateTime.now();
    final day = await showDatePicker(
      context: context,
      initialDate: _day,
      firstDate: DateUtils.dateOnly(now),
      lastDate: now.add(const Duration(days: 365)),
    );
    if (day != null) {
      setState(() => _day = day);
    }
  }

  Future<void> _pickTime() async {
    final time = await showTimePicker(context: context, initialTime: _time);
    if (time != null) {
      setState(() => _time = time);
    }
  }

  void _done() {
    final minutes = int.tryParse(_minutes.text.trim());
    final count = int.tryParse(_count.text.trim());
    final capacity = int.tryParse(_capacity.text.trim());
    if (minutes == null ||
        minutes < 5 ||
        minutes > 240 ||
        count == null ||
        count < 1 ||
        count > 50 ||
        capacity == null ||
        capacity < 1 ||
        capacity > 50) {
      setState(() => _invalid = true);
      return;
    }
    Navigator.of(context).pop(
      NewSlots(
        startsAt: DateTime(
          _day.year,
          _day.month,
          _day.day,
          _time.hour,
          _time.minute,
        ),
        minutes: minutes,
        place: _place.text.trim(),
        capacity: capacity,
        count: count,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    return AlertDialog(
      title: Text(l10n.defenseAddSlots),
      content: SizedBox(
        width: 420,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: _pickDay,
                      icon: const Icon(Icons.calendar_today_outlined),
                      label: Text(labels.day(_day)),
                    ),
                  ),
                  const SizedBox(width: Spacing.s),
                  OutlinedButton.icon(
                    onPressed: _pickTime,
                    icon: const Icon(Icons.schedule),
                    label: Text(_time.format(context)),
                  ),
                ],
              ),
              const SizedBox(height: Spacing.m),
              Row(
                children: [
                  Expanded(
                    child: TextField(
                      key: const ValueKey('slots-minutes'),
                      controller: _minutes,
                      keyboardType: TextInputType.number,
                      decoration: InputDecoration(
                        labelText: l10n.defenseSlotMinutes,
                      ),
                    ),
                  ),
                  const SizedBox(width: Spacing.m),
                  Expanded(
                    child: TextField(
                      key: const ValueKey('slots-count'),
                      controller: _count,
                      keyboardType: TextInputType.number,
                      decoration: InputDecoration(
                        labelText: l10n.defenseSlotCount,
                      ),
                    ),
                  ),
                  const SizedBox(width: Spacing.m),
                  Expanded(
                    child: TextField(
                      key: const ValueKey('slots-capacity'),
                      controller: _capacity,
                      keyboardType: TextInputType.number,
                      decoration: InputDecoration(
                        labelText: l10n.defenseSlotCapacity,
                      ),
                    ),
                  ),
                ],
              ),
              TextField(
                key: const ValueKey('slots-place'),
                controller: _place,
                maxLength: 500,
                decoration: InputDecoration(labelText: l10n.defenseSlotPlace),
              ),
              if (_invalid)
                Text(
                  l10n.defenseSlotsInvalid,
                  style: TextStyle(color: Theme.of(context).colorScheme.error),
                ),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(onPressed: _done, child: Text(l10n.defenseAdd)),
      ],
    );
  }
}
