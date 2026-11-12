class Loaded<T> {
  const Loaded(this.value, {required this.stale});

  final T value;
  final bool stale;
}
