import '../study/course_person.dart';
import 'leave_status.dart';

class LeaveRequest {
  const LeaveRequest({
    required this.id,
    required this.teamId,
    required this.teamName,
    required this.student,
    required this.reason,
    required this.status,
  });

  factory LeaveRequest.fromJson(Map<String, Object?> json) => LeaveRequest(
    id: json['id']! as String,
    teamId: json['teamId']! as String,
    teamName: json['teamName'] as String? ?? '',
    student: CoursePerson.fromJson(json['student']! as Map<String, Object?>),
    reason: json['reason'] as String? ?? '',
    status: LeaveStatus.parse(json['status']),
  );

  final String id;
  final String teamId;
  final String teamName;
  final CoursePerson student;
  final String reason;
  final LeaveStatus status;
}
