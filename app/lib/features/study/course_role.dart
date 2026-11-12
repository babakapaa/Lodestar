enum CourseRole {
  teacher,
  student,
  manager;

  static CourseRole parse(Object? value) => switch (value) {
    'teacher' => CourseRole.teacher,
    'manager' => CourseRole.manager,
    _ => CourseRole.student,
  };

  bool get teaches => this != CourseRole.student;
}
