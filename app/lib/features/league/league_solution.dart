import '../team/team_work.dart';
import 'similar_solution.dart';

class LeagueSolution {
  const LeagueSolution({
    required this.id,
    required this.teamName,
    required this.members,
    required this.work,
    required this.similarity,
    required this.explanation,
    required this.points,
    required this.passed,
    required this.comment,
    required this.reviewed,
    required this.similar,
  });

  factory LeagueSolution.fromJson(Map<String, Object?> json) {
    final similar = json['similar'];
    return LeagueSolution(
      id: json['id']! as String,
      teamName: json['teamName'] as String? ?? '',
      members: _ids(json['members']),
      work: TeamWork.fromJson(json['work']),
      similarity: (json['similarity'] as num?)?.toDouble(),
      explanation: json['explanation'] as String? ?? '',
      points: (json['points'] as num?)?.toInt(),
      passed: _ids(json['passed']),
      comment: json['comment'] as String? ?? '',
      reviewed: json['reviewedAt'] != null,
      similar: similar is Map<String, Object?>
          ? SimilarSolution.fromJson(similar)
          : null,
    );
  }

  static List<String> _ids(Object? value) => [
    for (final id in value as List<Object?>? ?? const []) id! as String,
  ];

  final String id;
  final String teamName;
  final List<String> members;
  final TeamWork work;
  final double? similarity;
  final String explanation;
  final int? points;
  final List<String> passed;
  final String comment;
  final bool reviewed;
  final SimilarSolution? similar;
}
