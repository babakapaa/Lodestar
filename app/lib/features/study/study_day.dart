abstract final class StudyDay {
  static String format(DateTime day) =>
      '${day.year.toString().padLeft(4, '0')}-'
      '${day.month.toString().padLeft(2, '0')}-'
      '${day.day.toString().padLeft(2, '0')}';

  static DateTime parse(String value) {
    final parsed = DateTime.tryParse(value);
    final day = parsed ?? DateTime.now();
    return DateTime(day.year, day.month, day.day);
  }

  static DateTime weekStart(DateTime day) =>
      DateTime(day.year, day.month, day.day - (day.weekday - 1));
}
