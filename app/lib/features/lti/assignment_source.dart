enum AssignmentSource {
  lodestar,
  lms;

  static AssignmentSource parse(Object? value) =>
      value == lms.name ? lms : lodestar;
}
