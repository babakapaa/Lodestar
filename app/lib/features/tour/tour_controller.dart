import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/settings/settings_providers.dart';

class TourController extends Notifier<bool> {
  static const key = 'tour.seen';

  @override
  bool build() => ref.watch(sharedPreferencesProvider).getBool(key) ?? false;

  Future<void> finish() async {
    state = true;
    await ref.read(sharedPreferencesProvider).setBool(key, true);
  }

  Future<void> restart() async {
    state = false;
    await ref.read(sharedPreferencesProvider).setBool(key, false);
  }
}
