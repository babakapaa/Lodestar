class SandboxTerm {
  const SandboxTerm({
    required this.id,
    required this.name,
    required this.startsOn,
    required this.endsOn,
  });

  factory SandboxTerm.fromJson(Map<String, Object?> json) => SandboxTerm(
    id: json['id']! as String,
    name: json['name'] as String? ?? '',
    startsOn: json['startsOn'] as String? ?? '',
    endsOn: json['endsOn'] as String? ?? '',
  );

  final String id;
  final String name;
  final String startsOn;
  final String endsOn;

  Map<String, Object?> toJson() => {
    'name': name,
    'startsOn': startsOn,
    'endsOn': endsOn,
  };
}
