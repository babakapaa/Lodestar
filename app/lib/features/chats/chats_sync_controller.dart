import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/api_exception.dart';
import '../../core/auth/auth_providers.dart';
import '../../core/auth/auth_state.dart';
import '../../core/live/live_event.dart';
import '../../core/live/live_providers.dart';
import 'chats_providers.dart';
import 'chats_sync.dart';
import 'chats_sync_status.dart';

class ChatsSyncController extends Notifier<ChatsSyncStatus> {
  static const events = {'chat.message', 'chat.changed', 'chat.read'};
  static const _debounce = Duration(milliseconds: 150);
  static const _firstDelay = Duration(seconds: 5);
  static const _maxDelay = Duration(seconds: 60);

  Timer? _timer;
  Duration _delay = _firstDelay;
  bool _running = false;
  bool _again = false;

  @override
  ChatsSyncStatus build() {
    ref.onDispose(() => _timer?.cancel());
    final local = ref
        .watch(chatsStoreProvider)
        .localChanges
        .listen((_) => _soon());
    ref.onDispose(local.cancel);
    ref.listen(liveEventsProvider, (_, next) {
      final type = next.value?.type;
      if (events.contains(type) || type == LiveEvent.resync) {
        _soon();
      }
    });
    if (ref.read(authProvider) is SignedIn) {
      Future.microtask(sync);
    }
    return ChatsSyncStatus.idle;
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
    state = ChatsSyncStatus.syncing;
    try {
      await ref.read(chatsSyncProvider).run();
      _delay = _firstDelay;
      state = ChatsSyncStatus.idle;
    } on ApiException catch (error) {
      if (!ChatsSync.offline(error)) {
        rethrow;
      }
      state = ChatsSyncStatus.offline;
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
