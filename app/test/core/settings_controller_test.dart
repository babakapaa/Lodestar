import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/settings/app_language.dart';
import 'package:lodestar/core/settings/appearance_mode.dart';
import 'package:lodestar/core/settings/settings_providers.dart';
import 'package:lodestar/design/tokens/theme_family.dart';
import 'package:shared_preferences/shared_preferences.dart';

Future<ProviderContainer> containerWith(Map<String, Object> stored) async {
  SharedPreferences.setMockInitialValues(stored);
  final prefs = await SharedPreferences.getInstance();
  return ProviderContainer(
    overrides: [sharedPreferencesProvider.overrideWithValue(prefs)],
  );
}

void main() {
  test('defaults when nothing is stored', () async {
    final container = await containerWith({});
    addTearDown(container.dispose);

    final settings = container.read(settingsProvider);

    expect(settings.family, ThemeFamily.campus);
    expect(settings.appearance, AppearanceMode.system);
    expect(settings.language, AppLanguage.system);
  });

  test('changes are saved and survive a restart', () async {
    final first = await containerWith({});
    final controller = first.read(settingsProvider.notifier);
    await controller.setFamily(ThemeFamily.lab);
    await controller.setAppearance(AppearanceMode.dark);
    await controller.setLanguage(AppLanguage.kk);
    first.dispose();

    final prefs = await SharedPreferences.getInstance();
    final second = ProviderContainer(
      overrides: [sharedPreferencesProvider.overrideWithValue(prefs)],
    );
    addTearDown(second.dispose);
    final settings = second.read(settingsProvider);

    expect(settings.family, ThemeFamily.lab);
    expect(settings.appearance, AppearanceMode.dark);
    expect(settings.language, AppLanguage.kk);
  });

  test('unknown stored values fall back to defaults', () async {
    final container = await containerWith({'settings.family': 'neon'});
    addTearDown(container.dispose);

    expect(container.read(settingsProvider).family, ThemeFamily.campus);
  });
}
