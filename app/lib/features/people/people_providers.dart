import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/auth_providers.dart';
import 'people_repository.dart';
import 'person_hit.dart';

final peopleRepositoryProvider = Provider<PeopleRepository>(
  (ref) => PeopleRepository(ref.watch(apiClientProvider)),
);

final peopleSearchProvider = FutureProvider.autoDispose
    .family<List<PersonHit>, String>(
      (ref, query) => ref.watch(peopleRepositoryProvider).search(query),
    );

final personNameProvider = FutureProvider.autoDispose.family<String, String>(
  (ref, userId) => ref.watch(peopleRepositoryProvider).name(userId),
);
