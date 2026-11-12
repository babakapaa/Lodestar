import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/api_exception.dart';
import '../../core/auth/auth_providers.dart';
import '../../core/auth/auth_state.dart';
import '../../core/live/live_event.dart';
import '../../core/live/live_providers.dart';
import 'notes_providers.dart';
import 'notes_sync.dart';
import 'notes_sync_status.dart';

class NotesSyncController extends Notifier<NotesSyncStatus> {
  static const changedEvent = 'notes.changed';
  static const _debounce = Duration(seconds: 1);
  static const _firstDelay = Duration(seconds: 5);
  static const _maxDelay = Duration(seconds: 60);

  Timer? _timer;
  Duration _delay = _firstDelay;
  bool _running = false;
  bool _again = false;

  @override
  NotesSyncStatus build() {
    ref.onDispose(() => _timer?.cancel());
    final local = ref
        .watch(notesStoreProvider)
        .localChanges
        .listen((_) => _soon());
    ref.onDispose(local.cancel);
    ref.listen(liveEventsProvider, (_, next) {
      final type = next.value?.type;
      if (type == changedEvent || type == LiveEvent.resync) {
        _soon();
      }
    });
    if (ref.read(authProvider) is SignedIn) {
      Future.microtask(sync);
    }
    return NotesSyncStatus.idle;
  }

  Future<void> sync() async {
    if (_running) {
      _again = true;
      return;
    }
    if (ref.read(authProvider) is! SignedIn) {
      return;
    }
    _running = true;
    _timer?.cancel();
    state = NotesSyncStatus.syncing;
    try {
      await ref.read(notesSyncProvider).run();
      _delay = _firstDelay;
      state = NotesSyncStatus.idle;
    } on ApiException catch (error) {
      if (!NotesSync.offline(error)) {
        rethrow;
      }
      state = NotesSyncStatus.offline;
      _retry();
    } finally {
      _running = false;
    }
    if (_again) {
      _again = false;
      await sync();
    }
  }

  void _soon() {
    _timer?.cancel();
    _timer = Timer(_debounce, sync);
  }

  void _retry() {
    _timer = Timer(_delay, sync);
    final next = _delay * 2;
    _delay = next > _maxDelay ? _maxDelay : next;
  }
}
