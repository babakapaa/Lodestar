class DeviceSession {
  const DeviceSession({
    required this.id,
    required this.deviceName,
    required this.lastUsedAt,
    required this.current,
  });

  factory DeviceSession.fromJson(Map<String, Object?> json) {
    return DeviceSession(
      id: json['id']! as String,
      deviceName: json['deviceName']! as String,
      lastUsedAt: DateTime.parse(json['lastUsedAt']! as String),
      current: json['current']! as bool,
    );
  }

  final String id;
  final String deviceName;
  final DateTime lastUsedAt;
  final bool current;
}
