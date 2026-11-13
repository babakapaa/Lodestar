class CallConsent {
  const CallConsent({
    required this.requestedBy,
    required this.audioOnly,
    required this.waiting,
  });

  factory CallConsent.fromJson(Map<String, Object?> json) => CallConsent(
    requestedBy: json['requestedBy']! as String,
    audioOnly: json['audioOnly'] as bool? ?? false,
    waiting: [
      for (final item in json['waiting'] as List<Object?>? ?? const [])
        item! as String,
    ],
  );

  final String requestedBy;
  final bool audioOnly;
  final List<String> waiting;
}
