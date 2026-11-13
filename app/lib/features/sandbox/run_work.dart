class RunWork {
  const RunWork({required this.kind, required this.id, required this.title});

  static RunWork? parse(Map<String, String> query) {
    final kind = query['workKind'];
    final id = query['workId'];
    if (kind == null || id == null) {
      return null;
    }
    return RunWork(kind: kind, id: id, title: query['workTitle'] ?? '');
  }

  final String kind;
  final String id;
  final String title;
}
