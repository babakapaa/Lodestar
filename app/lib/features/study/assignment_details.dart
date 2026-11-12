import 'assignment.dart';
import 'attached_file.dart';
import 'course_role.dart';
import 'student_work.dart';
import 'work_counts.dart';

class AssignmentDetails {
  const AssignmentDetails({
    required this.assignment,
    required this.courseTitle,
    required this.role,
    required this.files,
    required this.mine,
    required this.counts,
  });

  factory AssignmentDetails.fromJson(Map<String, Object?> json) {
    return AssignmentDetails(
      assignment: Assignment.fromJson(json),
      courseTitle: json['courseTitle'] as String? ?? '',
      role: CourseRole.parse(json['role']),
      files: AttachedFile.list(json['files']),
      mine: StudentWork.maybe(json['mine']),
      counts: WorkCounts.maybe(json['counts']),
    );
  }

  final Assignment assignment;
  final String courseTitle;
  final CourseRole role;
  final List<AttachedFile> files;
  final StudentWork? mine;
  final WorkCounts? counts;
}
