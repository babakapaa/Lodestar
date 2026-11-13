enum HistoryEvent {
  created,
  joined,
  added,
  left,
  removed,
  separated,
  movedIn('moved_in'),
  movedOut('moved_out');

  const HistoryEvent([this._wire]);

  final String? _wire;

  String get wire => _wire ?? name;

  static HistoryEvent parse(Object? value) => HistoryEvent.values.firstWhere(
    (e) => e.wire == value,
    orElse: () => HistoryEvent.joined,
  );
}
