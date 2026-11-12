import 'audit_entry.dart';

class AuditPage {
  const AuditPage({required this.entries, required this.hasMore});

  static const size = 50;

  final List<AuditEntry> entries;
  final bool hasMore;
}
