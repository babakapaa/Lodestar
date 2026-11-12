import 'course.dart';
import 'course_person.dart';

class CourseDetails {
  const CourseDetails({
    required this.course,
    required this.canManage,
    required this.teachers,
    required this.units,
  });

  factory CourseDetails.fromJson(Map<String, Object?> json) {
    return CourseDetails(
      course: Course.fromJson(json),
      canManage: json['canManage'] as bool? ?? false,
      teachers: [
        for (final item in json['teachers']! as List<Object?>)
          CoursePerson.fromJson(item! as Map<String, Object?>),
      ],
      units: [
        for (final item in json['units']! as List<Object?>)
          (item! as Map<String, Object?>)['name']! as String,
      ],
    );
  }

  final Course course;
  final bool canManage;
  final List<CoursePerson> teachers;
  final List<String> units;
}
