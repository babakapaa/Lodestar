import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/empty_state.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'defense_providers.dart';
import 'defense_slot.dart';
import 'slot_text.dart';

class DefenseBookingScreen extends ConsumerWidget {
  const DefenseBookingScreen({required this.assignmentId, super.key});

  final String assignmentId;

  Future<void> _book(
    BuildContext context,
    WidgetRef ref,
    DefenseSlot slot,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    final router = GoRouter.of(context);
    try {
      await ref.read(defenseApiProvider).book(assignmentId, slot.id);
      ref.invalidate(myDefenseProvider(assignmentId));
      router.pop();
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
      ref.invalidate(defenseSlotsProvider(assignmentId));
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final locale = Localizations.localeOf(context).toLanguageTag();
    final slots = ref.watch(defenseSlotsProvider(assignmentId));
    final now = DateTime.now();
    final list = [
      for (final slot in slots.value ?? const <DefenseSlot>[])
        if (slot.startsAt.isAfter(now)) slot,
    ];
    return Scaffold(
      appBar: AppBar(title: Text(l10n.defenseBookTitle)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              if (slots.hasError)
                FormError(message: ApiErrorText.of(l10n, slots.error!))
              else if (slots.value == null)
                const Center(child: CircularProgressIndicator())
              else if (list.isEmpty)
                EmptyState(
                  icon: Icons.event_busy_outlined,
                  title: l10n.defenseNoSlotsTitle,
                  body: l10n.defenseNoSlotsBody,
                )
              else
                for (final slot in list)
                  ListTile(
                    key: ValueKey('slot-${slot.id}'),
                    contentPadding: EdgeInsets.zero,
                    leading: Icon(
                      slot.mine ? Icons.event_available : Icons.event_outlined,
                      color: slot.mine ? context.palette.accent : null,
                    ),
                    title: Text(SlotText.of(slot, locale)),
                    subtitle: Text(
                      slot.mine
                          ? l10n.defenseSlotMine
                          : l10n.defenseSlotFree(slot.capacity - slot.taken),
                    ),
                    enabled: slot.mine || slot.free,
                    onTap: slot.mine ? null : () => _book(context, ref, slot),
                  ),
            ],
          ),
        ),
      ),
    );
  }
}
