import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'admin_providers.dart';
import 'setting_group.dart';
import 'setting_labels.dart';
import 'setting_tile.dart';
import '../../design/components/tile_group.dart';

class SettingsSection extends ConsumerWidget {
  const SettingsSection({super.key});

  Future<void> _save(
    BuildContext context,
    WidgetRef ref,
    String key,
    Object value,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await ref.read(adminRepositoryProvider).set(key, value);
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
    ref.invalidate(tenantSettingsProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final schema = ref.watch(settingsSchemaProvider);
    final values = ref.watch(tenantSettingsProvider);
    final error = schema.error ?? values.error;
    final definitions = schema.value;
    final current = values.value;
    final labels = SettingLabels(l10n);
    return Section(
      title: l10n.adminSettingsTitle,
      child: error != null
          ? FormError(message: ApiErrorText.of(l10n, error))
          : definitions == null || current == null
          ? const Center(child: CircularProgressIndicator())
          : Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                for (final group in SettingGroup.values)
                  if (definitions.any(
                    (d) => SettingGroup.of(d.key) == group,
                  )) ...[
                    Padding(
                      padding: const EdgeInsets.only(
                        top: Spacing.l,
                        bottom: Spacing.s,
                      ),
                      child: Text(
                        labels.group(group),
                        style: Theme.of(context).textTheme.titleSmall,
                      ),
                    ),
                    TileGroup(
                      children: [
                        for (final definition in definitions)
                          if (SettingGroup.of(definition.key) == group)
                            SettingTile(
                              definition: definition,
                              value: current[definition.key],
                              onChanged: (value) =>
                                  _save(context, ref, definition.key, value),
                            ),
                      ],
                    ),
                  ],
              ],
            ),
    );
  }
}
