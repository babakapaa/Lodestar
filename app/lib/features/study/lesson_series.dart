class LessonSeries {
  const LessonSeries({
    required this.kind,
    required this.room,
    required this.startsOn,
    required this.endsOn,
    required this.weekday,
    required this.startTime,
    required this.minutes,
    required this.everyWeeks,
    this.id = '',
  });

  factory LessonSeries.fromJson(Map<String, Object?> json) => LessonSeries(
    id: json['id']! as String,
    kind: json['kind'] as String? ?? '',
    room: json['room'] as String? ?? '',
    startsOn: json['startsOn']! as String,
    endsOn: json['endsOn']! as String,
    weekday: (json['weekday']! as num).toInt(),
    startTime: json['startTime']! as String,
    minutes: (json['minutes']! as num).toInt(),
    everyWeeks: (json['everyWeeks'] as num?)?.toInt() ?? 1,
  );

  final String id;
  final String kind;
  final String room;
  final String startsOn;
  final String endsOn;
  final int weekday;
  final String startTime;
  final int minutes;
  final int everyWeeks;

  Map<String, Object?> toJson() => {
    'kind': kind,
    'room': room,
    'startsOn': startsOn,
    'endsOn': endsOn,
    'weekday': weekday,
    'startTime': startTime,
    'minutes': minutes,
    'everyWeeks': everyWeeks,
  };
}
