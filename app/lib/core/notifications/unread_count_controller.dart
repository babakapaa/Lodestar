import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../live/live_event.dart';
import '../live/live_providers.dart';
import 'notifications_providers.dart';

class UnreadCountController extends AsyncNotifier<int> {
  @override
  Future<int> build() async {
    if (ref.watch(liveConnectionProvider) == null) {
      return 0;
    }
    ref.listen(liveEventsProvider, (previous, next) {
      if (next case AsyncData(:final value)) {
        if (value.type == LiveEvent.resync) {
          ref.invalidateSelf();
        } else if (value.type == LiveEvent.notification) {
          state = AsyncData((state.value ?? 0) + 1);
        }
      }
    });
    return ref.read(notificationsRepositoryProvider).unread();
  }
}
