import 'dart:ui';

enum AppLanguage {
  system(null),
  kk(Locale('kk')),
  ru(Locale('ru')),
  en(Locale('en'));

  const AppLanguage(this.locale);

  final Locale? locale;
}
