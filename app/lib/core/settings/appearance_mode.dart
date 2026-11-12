import 'package:material_ui/material_ui.dart';

enum AppearanceMode {
  system(ThemeMode.system),
  light(ThemeMode.light),
  dark(ThemeMode.dark);

  const AppearanceMode(this.themeMode);

  final ThemeMode themeMode;
}
