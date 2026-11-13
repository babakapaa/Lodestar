enum AttendanceStatus {
  present,
  late,
  absent,
  excused;

  static AttendanceStatus? parse(Object? value) {
    for (final status in values) {
      if (status.name == value) {
        return status;
      }
    }
    return null;
  }
}
