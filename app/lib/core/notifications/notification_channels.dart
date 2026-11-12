class NotificationChannels {
  const NotificationChannels({
    required this.type,
    required this.required,
    required this.inApp,
    required this.push,
    required this.email,
  });

  factory NotificationChannels.fromJson(Map<String, Object?> json) {
    return NotificationChannels(
      type: json['type']! as String,
      required: json['required']! as bool,
      inApp: json['inApp']! as bool,
      push: json['push']! as bool,
      email: json['email']! as bool,
    );
  }

  final String type;
  final bool required;
  final bool inApp;
  final bool push;
  final bool email;

  NotificationChannels copyWith({bool? inApp, bool? push, bool? email}) {
    return NotificationChannels(
      type: type,
      required: required,
      inApp: inApp ?? this.inApp,
      push: push ?? this.push,
      email: email ?? this.email,
    );
  }

  Map<String, Object?> toInput() => {
    'inApp': inApp,
    'push': push,
    'email': email,
  };
}
