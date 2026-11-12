import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../api/api_exception.dart';
import '../auth/auth_providers.dart';
import '../auth/auth_state.dart';
import 'live_connection.dart';
import 'live_event.dart';
import 'live_socket.dart';
import 'web_live_socket.dart';

final liveSocketProvider = Provider<LiveSocket>((ref) => const WebLiveSocket());

final liveConnectionProvider = Provider<LiveConnection?>((ref) {
  final auth = ref.watch(authProvider);
  if (auth is! SignedIn) {
    return null;
  }
  final connection = LiveConnection(
    socket: ref.watch(liveSocketProvider),
    address: auth.address,
    tokens: ref.watch(tokenStoreProvider),
    onFailure: () async {
      try {
        await ref.read(apiClientProvider).get('/v1/notifications/unread');
      } on ApiException {
        return;
      }
    },
  )..start();
  ref.onDispose(connection.close);
  return connection;
});

final liveEventsProvider = StreamProvider<LiveEvent>((ref) {
  final connection = ref.watch(liveConnectionProvider);
  return connection?.events ?? const Stream.empty();
});
