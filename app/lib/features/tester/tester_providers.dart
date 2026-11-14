import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'tester_api.dart';
import 'tester_info.dart';
import 'tester_next.dart';
import 'tester_person.dart';
import 'tester_switch.dart';

final testerApiProvider = Provider<TesterApi>(
  (ref) => TesterApi(ref.watch(apiTransportProvider)),
);

final testerInfoProvider = FutureProvider.autoDispose<TesterInfo?>((ref) {
  final address = ref.watch(authProvider.select((state) => state.address));
  if (address == null) {
    return null;
  }
  return ref.watch(testerApiProvider).info(address);
});

final testerSwitchProvider = Provider<TesterSwitch>(TesterSwitch.new);

final testerNextProvider = NotifierProvider<TesterNext, TesterPerson?>(
  TesterNext.new,
);
