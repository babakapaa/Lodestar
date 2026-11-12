import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'admin_providers.dart';
import 'audit_page.dart';

class AuditController extends AsyncNotifier<AuditPage> {
  @override
  Future<AuditPage> build() async {
    final entries = await ref.read(adminRepositoryProvider).audit();
    return AuditPage(
      entries: entries,
      hasMore: entries.length >= AuditPage.size,
    );
  }

  Future<void> loadMore() async {
    final current = state.value;
    if (current == null || !current.hasMore || current.entries.isEmpty) {
      return;
    }
    final next = await ref
        .read(adminRepositoryProvider)
        .audit(before: current.entries.last.id);
    state = AsyncData(
      AuditPage(
        entries: [...current.entries, ...next],
        hasMore: next.length >= AuditPage.size,
      ),
    );
  }
}
