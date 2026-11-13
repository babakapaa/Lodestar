import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'defense.dart';
import 'defense_api.dart';
import 'defense_row.dart';
import 'defense_slot.dart';
import 'my_defense.dart';

final defenseApiProvider = Provider<DefenseApi>(
  (ref) => DefenseApi(ref.watch(apiClientProvider)),
);

final myDefenseProvider = FutureProvider.autoDispose.family<MyDefense, String>(
  (ref, id) => ref.watch(defenseApiProvider).mine(id),
);

final defenseSlotsProvider = FutureProvider.autoDispose
    .family<List<DefenseSlot>, String>(
      (ref, id) => ref.watch(defenseApiProvider).slots(id),
    );

final defenseRowsProvider = FutureProvider.autoDispose
    .family<List<DefenseRow>, String>(
      (ref, id) => ref.watch(defenseApiProvider).list(id),
    );

final defenseProvider = FutureProvider.autoDispose
    .family<Defense, (String, String)>(
      (ref, key) => ref.watch(defenseApiProvider).get(key.$1, key.$2),
    );
