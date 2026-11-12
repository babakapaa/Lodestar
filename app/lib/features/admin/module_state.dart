class ModuleState {
  const ModuleState({
    required this.name,
    required this.optional,
    required this.enabled,
  });

  factory ModuleState.fromJson(Map<String, Object?> json) {
    return ModuleState(
      name: json['name']! as String,
      optional: json['optional']! as bool,
      enabled: json['enabled']! as bool,
    );
  }

  final String name;
  final bool optional;
  final bool enabled;
}
