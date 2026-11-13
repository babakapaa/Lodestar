class SetupStatus {
  const SetupStatus({
    required this.launched,
    required this.manages,
    required this.assigned,
  });

  factory SetupStatus.fromJson(Map<String, Object?> json) => SetupStatus(
    launched: json['launched'] as bool? ?? true,
    manages: json['manages'] as bool? ?? false,
    assigned: [
      for (final item in json['assigned'] as List<Object?>? ?? const [])
        item! as String,
    ],
  );

  static const open = SetupStatus(launched: true, manages: false, assigned: []);

  final bool launched;
  final bool manages;
  final List<String> assigned;

  bool get takesPart => !launched && (manages || assigned.isNotEmpty);
}
