import 'search_map.dart';

class SubmissionDraft {
  const SubmissionDraft({
    required this.clientId,
    required this.assignmentId,
    required this.assignmentTitle,
    required this.body,
    required this.composedAt,
    this.searchMap,
    this.error,
  });

  factory SubmissionDraft.fromJson(Map<String, Object?> json) {
    return SubmissionDraft(
      clientId: json['clientId']! as String,
      assignmentId: json['assignmentId']! as String,
      assignmentTitle: json['assignmentTitle'] as String? ?? '',
      body: json['body'] as String? ?? '',
      composedAt: DateTime.parse(json['composedAt']! as String),
      searchMap: SearchMap.maybe(json['searchMap']),
      error: json['error'] as String?,
    );
  }

  final String clientId;
  final String assignmentId;
  final String assignmentTitle;
  final String body;
  final DateTime composedAt;
  final SearchMap? searchMap;
  final String? error;

  bool get failed => error != null;

  SubmissionDraft withError(String? code) => SubmissionDraft(
    clientId: clientId,
    assignmentId: assignmentId,
    assignmentTitle: assignmentTitle,
    body: body,
    composedAt: composedAt,
    searchMap: searchMap,
    error: code,
  );

  Map<String, Object?> toJson() => {
    'clientId': clientId,
    'assignmentId': assignmentId,
    'assignmentTitle': assignmentTitle,
    'body': body,
    'composedAt': composedAt.toUtc().toIso8601String(),
    'searchMap': searchMap?.toJson(),
    'error': ?error,
  };

  Map<String, Object?> toRequest() => {
    'clientId': clientId,
    'body': body,
    'composedAt': composedAt.toUtc().toIso8601String(),
    'searchMap': searchMap?.toJson(),
  };
}
