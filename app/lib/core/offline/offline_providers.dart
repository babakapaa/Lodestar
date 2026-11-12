import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../auth/auth_providers.dart';
import '../settings/settings_providers.dart';
import 'cached_reader.dart';
import 'response_cache.dart';

final responseCacheProvider = Provider<ResponseCache>(
  (ref) => ResponseCache(ref.watch(sharedPreferencesProvider)),
);

final cachedReaderProvider = Provider<CachedReader>(
  (ref) => CachedReader(
    ref.watch(apiClientProvider),
    ref.watch(responseCacheProvider),
  ),
);
