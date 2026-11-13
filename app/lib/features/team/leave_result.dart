import 'leave_request.dart';

class LeaveResult {
  const LeaveResult({required this.left, required this.request});

  factory LeaveResult.fromJson(Map<String, Object?> json) {
    final request = json['request'];
    return LeaveResult(
      left: json['left'] as bool? ?? false,
      request: request is Map<String, Object?>
          ? LeaveRequest.fromJson(request)
          : null,
    );
  }

  final bool left;
  final LeaveRequest? request;
}
