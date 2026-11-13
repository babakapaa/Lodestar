import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../l10n/l10n_context.dart';
import '../admin/admin_providers.dart';
import '../admin/setting_tile.dart';
import '../team/team_action.dart';
import 'setup_providers.dart';
import 'setup_section.dart';

class SectionSettings extends ConsumerWidget {
  const SectionSettings({required this.section, super.key});

  final SetupSection section;

  Future<void> _save(
    BuildContext context,
    WidgetRef ref,
    String key,
    Object value,
  ) async {
    await TeamAction.run(
      context,
      () => ref.read(setupApiProvider).setSetting(section.key, key, value),
    );
    ref
      ..invalidate(tenantSettingsProvider)
      ..invalidate(setupOverviewProvider);
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    if (section.settings.isEmpty) {
      return Text(l10n.setupNoSettings);
    }
    final schema = ref.watch(settingsSchemaProvider);
    final values = ref.watch(tenantSettingsProvider);
    final error = schema.error ?? values.error;
    final definitions = schema.value;
    final current = values.value;
    if (error != null) {
      return FormError(message: ApiErrorText.of(l10n, error));
    }
    if (definitions == null || current == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return Column(
      children: [
        for (final definition in definitions)
          if (section.settings.contains(definition.key))
            SettingTile(
              definition: definition,
              value: current[definition.key],
              onChanged: (value) => _save(context, ref, definition.key, value),
            ),
      ],
    );
  }
}
