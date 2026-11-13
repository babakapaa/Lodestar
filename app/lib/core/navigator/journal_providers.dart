import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../auth/auth_providers.dart';
import 'journal_api.dart';
import 'journal_entry.dart';

final journalApiProvider = Provider<JournalApi>(
  (ref) => JournalApi(ref.watch(apiClientProvider)),
);

final journalProvider = FutureProvider.autoDispose<List<JournalEntry>>(
  (ref) => ref.watch(journalApiProvider).list(),
);

final assignmentJournalProvider = FutureProvider.autoDispose
    .family<List<JournalEntry>, String>(
      (ref, assignmentId) => ref
          .watch(journalApiProvider)
          .list(scopeKind: 'assignment', scopeId: assignmentId),
    );
