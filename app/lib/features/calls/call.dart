import 'call_consent.dart';

class Call {
  const Call({
    required this.id,
    required this.chatId,
    required this.startedBy,
    required this.startedAt,
    required this.people,
    required this.recording,
    this.audioOnly = false,
    this.consent,
  });

  factory Call.fromJson(Map<String, Object?> json) => Call(
    id: json['id']! as String,
    chatId: json['chatId']! as String,
    startedBy: json['startedBy'] as String?,
    startedAt: DateTime.parse(json['startedAt']! as String),
    people: [
      for (final item in json['people'] as List<Object?>? ?? const [])
        item! as String,
    ],
    recording: json['recording'] as bool? ?? false,
    audioOnly: json['audioOnly'] as bool? ?? false,
    consent: switch (json['consent']) {
      final Map<String, Object?> c => CallConsent.fromJson(c),
      _ => null,
    },
  );

  final String id;
  final String chatId;
  final String? startedBy;
  final DateTime startedAt;
  final List<String> people;
  final bool recording;
  final bool audioOnly;
  final CallConsent? consent;
}
