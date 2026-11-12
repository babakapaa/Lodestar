import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'contact_draft.dart';
import 'contact_kinds.dart';

class ContactRow extends StatelessWidget {
  const ContactRow({
    required this.draft,
    required this.onChanged,
    required this.onRemove,
    super.key,
  });

  final ContactDraft draft;
  final VoidCallback onChanged;
  final VoidCallback onRemove;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final kinds = {...ContactKinds.presets, draft.kind};
    return Padding(
      padding: const EdgeInsets.only(bottom: Spacing.m),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              SizedBox(
                width: 140,
                child: DropdownButtonFormField<String>(
                  initialValue: draft.kind,
                  isExpanded: true,
                  decoration: InputDecoration(labelText: l10n.contactKind),
                  items: [
                    for (final kind in kinds)
                      DropdownMenuItem(
                        value: kind,
                        child: Text(ContactKinds.label(l10n, kind)),
                      ),
                  ],
                  onChanged: (kind) {
                    if (kind != null) {
                      draft.kind = kind;
                      onChanged();
                    }
                  },
                ),
              ),
              const SizedBox(width: Spacing.m),
              Expanded(
                child: TextField(
                  controller: draft.value,
                  maxLength: 200,
                  decoration: InputDecoration(
                    labelText: l10n.contactValue,
                    counterText: '',
                  ),
                ),
              ),
              IconButton(
                tooltip: l10n.contactRemove,
                onPressed: onRemove,
                icon: const Icon(Icons.close),
              ),
            ],
          ),
          SwitchListTile(
            contentPadding: EdgeInsets.zero,
            title: Text(l10n.contactPublic),
            value: draft.public,
            onChanged: (public) {
              draft.public = public;
              onChanged();
            },
          ),
        ],
      ),
    );
  }
}
