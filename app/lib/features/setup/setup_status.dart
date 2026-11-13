class SetupStatus {
  const SetupStatus({
    required this.launched,
    required this.manages,
    required this.assigned,
    this.pilot = false,
    this.included = true,
  });

  factory SetupStatus.fromJson(Map<String, Object?> json) => SetupStatus(
    launched: json['launched'] as bool? ?? true,
    manages: json['manages'] as bool? ?? false,
    assigned: [
      for (final item in json['assigned'] as List<Object?>? ?? const [])
        item! as String,
    ],
    pilot: json['pilot'] as bool? ?? false,
    included: json['included'] as bool? ?? true,
  );

  static const open = SetupStatus(launched: true, manages: false, assigned: []);

  final bool launched;
  final bool manages;
  final List<String> assigned;
  final bool pilot;
  final bool included;

  bool get takesPart => !launched && (manages || assigned.isNotEmpty);

  bool get outsidePilot => pilot && !included;

  bool get opensWizard => manages && (takesPart || pilot);
}
