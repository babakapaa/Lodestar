import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/auth/auth_state.dart';
import '../chats/chats_providers.dart';
import '../notes/notes_providers.dart';
import 'tester_person.dart';
import 'tester_providers.dart';

class TesterSwitch {
  const TesterSwitch(this._ref);

  final Ref _ref;

  Future<void> signInAs(TesterPerson person) async {
    final auth = _ref.read(authProvider.notifier);
    final state = _ref.read(authProvider);
    final address = state.address;
    if (address == null) {
      return;
    }
    if (state is SignedIn) {
      _ref.read(testerNextProvider.notifier).set(person);
      await auth.signOut();
      return;
    }
    final tokens = await _ref
        .read(testerApiProvider)
        .signIn(address, person, device: _ref.read(deviceNameProvider));
    await _ref.read(notesStoreProvider).clear();
    await _ref.read(chatsStoreProvider).clear();
    await auth.adopt(tokens);
  }
}
