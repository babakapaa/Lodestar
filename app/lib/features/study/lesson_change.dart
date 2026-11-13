class LessonChange {
  const LessonChange.cancel({this.note = ''})
    : cancelled = true,
      startsAt = null,
      room = null;

  const LessonChange.move({this.startsAt, this.room, this.note = ''})
    : cancelled = false;

  final bool cancelled;
  final DateTime? startsAt;
  final String? room;
  final String note;

  Map<String, Object?> toJson() => {
    'status': cancelled ? 'cancelled' : 'changed',
    'startsAt': ?startsAt?.toUtc().toIso8601String(),
    'room': ?room,
    'note': note,
  };
}
