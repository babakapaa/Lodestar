import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/defense_format.dart';
import '../study/study_providers.dart';
import 'add_slots_dialog.dart';
import 'defense_providers.dart';
import 'defense_row.dart';
import 'defense_row_tile.dart';
import 'defense_slot.dart';
import 'sample_dialog.dart';
import 'slot_text.dart';

class DefensesScreen extends ConsumerWidget {
  const DefensesScreen({required this.assignmentId, super.key});

  final String assignmentId;

  Future<void> _run(
    BuildContext context,
    WidgetRef ref,
    Future<void> Function() request,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await request();
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref
      ..invalidate(defenseSlotsProvider(assignmentId))
      ..invalidate(defenseRowsProvider(assignmentId));
  }

  Future<void> _addSlots(BuildContext context, WidgetRef ref) async {
    final slots = await AddSlotsDialog.show(context);
    if (slots == null || !context.mounted) {
      return;
    }
    await _run(
      context,
      ref,
      () => ref.read(defenseApiProvider).addSlots(assignmentId, slots),
    );
  }

  Future<void> _sample(BuildContext context, WidgetRef ref) async {
    final count = await SampleDialog.show(context);
    if (count == null || !context.mounted) {
      return;
    }
    await _run(
      context,
      ref,
      () => ref.read(defenseApiProvider).sample(assignmentId, count),
    );
  }

  Widget _slot(BuildContext context, WidgetRef ref, DefenseSlot slot) {
    final l10n = context.l10n;
    final locale = Localizations.localeOf(context).toLanguageTag();
    return ListTile(
      key: ValueKey('slot-${slot.id}'),
      contentPadding: EdgeInsets.zero,
      leading: const Icon(Icons.event_outlined),
      title: Text(SlotText.of(slot, locale)),
      subtitle: Text(
        slot.students.isEmpty
            ? l10n.defenseSlotEmpty
            : [for (final s in slot.students) s.label].join(', '),
      ),
      trailing: IconButton(
        tooltip: l10n.defenseDeleteSlot,
        icon: const Icon(Icons.delete_outline),
        onPressed: () => _run(
          context,
          ref,
          () => ref.read(defenseApiProvider).deleteSlot(slot.id),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final assignment = ref
        .watch(assignmentProvider(assignmentId))
        .value
        ?.value
        .assignment;
    final slots = ref.watch(defenseSlotsProvider(assignmentId));
    final rows = ref.watch(defenseRowsProvider(assignmentId));
    final error = slots.error ?? rows.error;
    final sampled = assignment?.defense == DefenseFormat.sample;
    return Scaffold(
      appBar: AppBar(title: Text(assignment?.title ?? l10n.defensesTitle)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: error != null
              ? Padding(
                  padding: const EdgeInsets.all(Spacing.l),
                  child: FormError(message: ApiErrorText.of(l10n, error)),
                )
              : ListView(
                  padding: const EdgeInsets.all(Spacing.l),
                  children: [
                    Section(
                      title: l10n.defenseSlots,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          for (final slot
                              in slots.value ?? const <DefenseSlot>[])
                            _slot(context, ref, slot),
                          const SizedBox(height: Spacing.s),
                          Align(
                            alignment: Alignment.centerLeft,
                            child: OutlinedButton.icon(
                              key: const ValueKey('add-slots'),
                              onPressed: () => _addSlots(context, ref),
                              icon: const Icon(Icons.add),
                              label: Text(l10n.defenseAddSlots),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: Spacing.xl),
                    Section(
                      title: l10n.defenseStudents,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          if (sampled) ...[
                            Text(l10n.defenseSampleHint),
                            Align(
                              alignment: Alignment.centerLeft,
                              child: TextButton.icon(
                                key: const ValueKey('sample'),
                                onPressed: () => _sample(context, ref),
                                icon: const Icon(Icons.shuffle),
                                label: Text(l10n.defenseSample),
                              ),
                            ),
                          ],
                          if (rows.value == null)
                            const Center(child: CircularProgressIndicator())
                          else
                            for (final row
                                in rows.value ?? const <DefenseRow>[])
                              DefenseRowTile(
                                key: ValueKey(row.student.userId),
                                assignmentId: assignmentId,
                                row: row,
                                sampled: sampled,
                                onSelect: (selected) => _run(
                                  context,
                                  ref,
                                  () => ref
                                      .read(defenseApiProvider)
                                      .select(
                                        assignmentId,
                                        row.student.userId,
                                        selected,
                                      ),
                                ),
                              ),
                        ],
                      ),
                    ),
                  ],
                ),
        ),
      ),
    );
  }
}
