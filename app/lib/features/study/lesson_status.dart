enum LessonStatus {
  planned,
  changed,
  cancelled;

  static LessonStatus parse(Object? value) => switch (value) {
    'changed' => LessonStatus.changed,
    'cancelled' => LessonStatus.cancelled,
    _ => LessonStatus.planned,
  };
}
