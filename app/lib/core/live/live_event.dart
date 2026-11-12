class LiveEvent {
  const LiveEvent({
    required this.id,
    required this.type,
    required this.payload,
  });

  factory LiveEvent.fromJson(Map<String, Object?> json) {
    return LiveEvent(
      id: (json['id']! as num).toInt(),
      type: json['type']! as String,
      payload: json['payload'],
    );
  }

  static const notification = 'core.notification';
  static const resync = 'core.resync';

  final int id;
  final String type;
  final Object? payload;
}
