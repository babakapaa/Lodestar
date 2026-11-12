import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import '../people/people_providers.dart';
import 'audit_entry.dart';
import 'audit_labels.dart';

class AuditTile extends ConsumerWidget {
  const AuditTile({required this.entry, super.key});

  final AuditEntry entry;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final actorId = entry.actorId;
    final actor = actorId == null
        ? l10n.auditBySystem
        : ref.watch(personNameProvider(actorId)).value ?? '';
    final time = DateFormat.yMMMd(
      Localizations.localeOf(context).toLanguageTag(),
    ).add_Hm().format(entry.at.toLocal());
    final target = entry.details['name'] ?? entry.details['key'];
    return ListTile(
      leading: const Icon(Icons.history),
      title: Text(
        [
          AuditLabels(l10n).action(entry.action),
          if (target != null) '$target',
        ].join(': '),
      ),
      subtitle: Text([if (actor.isNotEmpty) actor, time].join(', ')),
    );
  }
}
