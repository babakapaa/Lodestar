import '../study/course_person.dart';

abstract final class TeamPeople {
  static List<CoursePerson> parse(Object? value) => [
    for (final item in value as List<Object?>? ?? const [])
      CoursePerson.fromJson(item! as Map<String, Object?>),
  ];
}
