enum RunStatus {
  queued,
  starting,
  running,
  held,
  stopping,
  finished,
  failed,
  stopped;

  static RunStatus parse(String? value) =>
      values.firstWhere((s) => s.name == value, orElse: () => failed);

  bool get active =>
      this == queued ||
      this == starting ||
      this == running ||
      this == held ||
      this == stopping;
}
