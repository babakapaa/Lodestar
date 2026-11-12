class CachedBody {
  const CachedBody(this.body, {required this.stale});

  final Object? body;
  final bool stale;
}
