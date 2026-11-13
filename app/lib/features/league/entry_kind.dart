enum EntryKind {
  work,
  adjustment,
  task;

  static EntryKind parse(Object? value) => EntryKind.values.firstWhere(
    (k) => k.name == value,
    orElse: () => EntryKind.work,
  );
}
