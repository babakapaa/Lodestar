import 'package:cupertino_ui/cupertino_ui.dart'
    show CupertinoPageTransitionsBuilder;
import 'package:material_ui/material_ui.dart';

import 'tokens/accent_choice.dart';
import 'tokens/lodestar_tokens.dart';
import 'tokens/theme_family.dart';
import 'typography/lodestar_typography.dart';

class LodestarTheme {
  const LodestarTheme(this.family, this.accent);

  final ThemeFamily family;
  final AccentChoice accent;

  ThemeData light() => _build(Brightness.light);

  ThemeData dark() => _build(Brightness.dark);

  ThemeData _build(Brightness brightness) {
    final spec = family.spec;
    final palette = accent.apply(spec.palette(brightness), brightness);
    final tokens = LodestarTokens.from(spec, palette);
    final text = LodestarTypography(spec, palette).build();
    final scheme = ColorScheme(
      brightness: brightness,
      primary: palette.accent,
      onPrimary: palette.onAccent,
      primaryContainer: palette.accentSoft,
      onPrimaryContainer: palette.accent,
      secondary: palette.accent,
      onSecondary: palette.onAccent,
      error: palette.warn,
      onError: palette.bg,
      surface: palette.surface,
      onSurface: palette.text,
      onSurfaceVariant: palette.muted,
      surfaceContainerLowest: palette.bg,
      surfaceContainerLow: palette.side,
      surfaceContainer: palette.surface,
      surfaceContainerHigh: palette.surface2,
      outline: palette.border,
      outlineVariant: palette.border,
    );
    final selectedLabel = text.labelMedium!.copyWith(color: palette.accent);
    final idleLabel = text.labelMedium!.copyWith(color: palette.muted);

    return ThemeData(
      useMaterial3: true,
      brightness: brightness,
      colorScheme: scheme,
      scaffoldBackgroundColor: palette.bg,
      canvasColor: palette.bg,
      dividerColor: palette.border,
      textTheme: text,
      splashFactory: InkSparkle.splashFactory,
      extensions: [tokens],
      dividerTheme: DividerThemeData(
        color: palette.border,
        thickness: 1,
        space: 1,
      ),
      iconTheme: IconThemeData(color: palette.text, size: 20),
      appBarTheme: AppBarTheme(
        backgroundColor: palette.bg,
        foregroundColor: palette.text,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        titleTextStyle: text.titleLarge,
        shape: Border(bottom: BorderSide(color: palette.border)),
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: palette.bg,
        indicatorColor: palette.accentSoft,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        height: 64,
        labelTextStyle: WidgetStateProperty.resolveWith(
          (states) =>
              states.contains(WidgetState.selected) ? selectedLabel : idleLabel,
        ),
        iconTheme: WidgetStateProperty.resolveWith(
          (states) => IconThemeData(
            color: states.contains(WidgetState.selected)
                ? palette.accent
                : palette.muted,
            size: 22,
          ),
        ),
      ),
      navigationRailTheme: NavigationRailThemeData(
        backgroundColor: palette.side,
        indicatorColor: palette.accentSoft,
        indicatorShape: RoundedRectangleBorder(
          borderRadius: tokens.borderRadiusSmall,
        ),
        selectedIconTheme: IconThemeData(color: palette.accent, size: 20),
        unselectedIconTheme: IconThemeData(color: palette.muted, size: 20),
        selectedLabelTextStyle: selectedLabel,
        unselectedLabelTextStyle: text.labelMedium,
      ),
      tooltipTheme: TooltipThemeData(
        decoration: BoxDecoration(
          color: palette.text,
          borderRadius: tokens.borderRadiusSmall,
        ),
        textStyle: text.bodySmall!.copyWith(color: palette.bg),
      ),
      pageTransitionsTheme: const PageTransitionsTheme(
        builders: {
          TargetPlatform.android: FadeForwardsPageTransitionsBuilder(),
          TargetPlatform.iOS: CupertinoPageTransitionsBuilder(),
          TargetPlatform.macOS: FadeForwardsPageTransitionsBuilder(),
          TargetPlatform.windows: FadeForwardsPageTransitionsBuilder(),
          TargetPlatform.linux: FadeForwardsPageTransitionsBuilder(),
        },
      ),
    );
  }
}
