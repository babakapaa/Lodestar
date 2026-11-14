import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/settings/app_language.dart';
import '../../core/settings/appearance_mode.dart';
import '../../core/settings/settings_providers.dart';
import '../../design/components/choice_group.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/theme_family.dart';
import '../../l10n/l10n_context.dart';
import 'accent_picker.dart';
import 'family_swatch.dart';
import 'settings_labels.dart';

class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final settings = ref.watch(settingsProvider);
    final controller = ref.read(settingsProvider.notifier);
    final labels = SettingsLabels(context.l10n);
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.settingsTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: [
                Section(
                  title: context.l10n.settingsTheme,
                  child: ChoiceGroup<ThemeFamily>(
                    values: ThemeFamily.values,
                    selected: settings.family,
                    labelOf: labels.family,
                    leadingOf: (family) => FamilySwatch(family: family),
                    onSelected: controller.setFamily,
                  ),
                ),
                const SizedBox(height: Spacing.xl),
                Section(
                  title: context.l10n.settingsAccent,
                  child: AccentPicker(
                    family: settings.family,
                    selected: settings.accent,
                    labelOf: labels.accent,
                    onSelected: controller.setAccent,
                  ),
                ),
                const SizedBox(height: Spacing.xl),
                Section(
                  title: context.l10n.settingsAppearance,
                  child: ChoiceGroup<AppearanceMode>(
                    values: AppearanceMode.values,
                    selected: settings.appearance,
                    labelOf: labels.appearance,
                    onSelected: controller.setAppearance,
                  ),
                ),
                const SizedBox(height: Spacing.xl),
                Section(
                  title: context.l10n.settingsLanguage,
                  child: ChoiceGroup<AppLanguage>(
                    values: AppLanguage.values,
                    selected: settings.language,
                    labelOf: labels.language,
                    onSelected: controller.setLanguage,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
