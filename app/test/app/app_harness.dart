import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/app/lodestar_app.dart';
import 'package:lodestar/core/settings/settings_providers.dart';
import 'package:material_ui/material_ui.dart';
import 'package:shared_preferences/shared_preferences.dart';

Future<void> pumpApp(
  WidgetTester tester, {
  required Size size,
  Map<String, Object> stored = const {'settings.language': 'ru'},
}) async {
  tester.view.physicalSize = size;
  tester.view.devicePixelRatio = 1;
  addTearDown(tester.view.reset);
  SharedPreferences.setMockInitialValues(stored);
  final prefs = await SharedPreferences.getInstance();
  await tester.pumpWidget(
    ProviderScope(
      overrides: [sharedPreferencesProvider.overrideWithValue(prefs)],
      child: const LodestarApp(),
    ),
  );
  await tester.pumpAndSettle();
}
