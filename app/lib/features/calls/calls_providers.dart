import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import '../../core/live/live_providers.dart';
import 'call.dart';
import 'call_engine.dart';
import 'calls_api.dart';
import 'livekit_call_engine.dart';

const callEvent = 'chat.call';

final callsApiProvider = Provider<CallsApi>(
  (ref) => CallsApi(ref.watch(apiClientProvider)),
);

final callEngineFactoryProvider = Provider<CallEngine Function()>(
  (ref) => LiveKitCallEngine.new,
);

final activeCallProvider = FutureProvider.autoDispose.family<Call?, String>((
  ref,
  chatId,
) {
  ref.listen(liveEventsProvider, (_, next) {
    final event = next.value;
    if (event != null &&
        event.type == callEvent &&
        event.payload is Map<String, Object?> &&
        (event.payload! as Map<String, Object?>)['chatId'] == chatId) {
      ref.invalidateSelf();
    }
  });
  return ref.watch(callsApiProvider).active(chatId);
});
