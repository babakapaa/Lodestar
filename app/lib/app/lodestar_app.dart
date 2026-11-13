import 'package:cupertino_ui/cupertino_ui.dart'
    show GlobalCupertinoLocalizations;
import 'package:flutter_localizations/flutter_localizations.dart'
    show GlobalWidgetsLocalizations;
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../core/auth/auth_providers.dart';
import '../core/auth/auth_state.dart';
import '../core/settings/settings_providers.dart';
import '../design/lodestar_theme.dart';
import '../features/chats/chats_providers.dart';
import '../features/notes/notes_providers.dart';
import '../design/tokens/motion.dart';
import '../l10n/generated/app_localizations.dart';
import 'router/app_router.dart';

class LodestarApp extends ConsumerStatefulWidget {
  const LodestarApp({super.key});

  @override
  ConsumerState<LodestarApp> createState() => _LodestarAppState();
}

class _LodestarAppState extends ConsumerState<LodestarApp> {
  late final ValueNotifier<AuthState> _auth;
  late final AppRouter _router;

  @override
  void initState() {
    super.initState();
    _auth = ValueNotifier(ref.read(authProvider));
    _router = AppRouter(_auth);
    ref.listenManual(authProvider, (previous, next) {
      _auth.value = next;
      if (next is SignedIn && previous is! SignedIn) {
        ref.read(settingsProvider.notifier).pull();
      }
      if (next is Disconnected &&
          previous != null &&
          previous is! Disconnected) {
        ref.read(notesStoreProvider).clear();
        ref.read(chatsStoreProvider).clear();
      }
    }, fireImmediately: true);
  }

  @override
  void dispose() {
    _router.config.dispose();
    _auth.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final settings = ref.watch(settingsProvider);
    final theme = LodestarTheme(settings.family);
    return MaterialApp.router(
      onGenerateTitle: (context) => AppLocalizations.of(context).appTitle,
      debugShowCheckedModeBanner: false,
      theme: theme.light(),
      darkTheme: theme.dark(),
      themeMode: settings.appearance.themeMode,
      themeAnimationDuration: Motion.slow,
      locale: settings.language.locale,
      supportedLocales: AppLocalizations.supportedLocales,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      routerConfig: _router.config,
    );
  }
}
