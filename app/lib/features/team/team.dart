import 'team_member.dart';

class Team {
  const Team({
    required this.id,
    required this.scopeKind,
    required this.scopeId,
    required this.name,
    required this.members,
    required this.archived,
  });

  factory Team.fromJson(Map<String, Object?> json) => Team(
    id: json['id']! as String,
    scopeKind: json['scopeKind']! as String,
    scopeId: json['scopeId']! as String,
    name: json['name'] as String? ?? '',
    members: [
      for (final item in json['members'] as List<Object?>? ?? const [])
        TeamMember.fromJson(item! as Map<String, Object?>),
    ],
    archived: json['archivedAt'] != null,
  );

  final String id;
  final String scopeKind;
  final String scopeId;
  final String name;
  final List<TeamMember> members;
  final bool archived;

  bool has(String userId) => members.any((m) => m.person.userId == userId);

  String get names => members.map((m) => m.person.label).join(', ');
}
