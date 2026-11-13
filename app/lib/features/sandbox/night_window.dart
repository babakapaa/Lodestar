class NightWindow {
  const NightWindow({
    required this.from,
    required this.to,
    required this.discount,
    required this.timeZone,
  });

  factory NightWindow.fromJson(Map<String, Object?> json) => NightWindow(
    from: json['from'] as int? ?? 0,
    to: json['to'] as int? ?? 0,
    discount: json['discount'] as int? ?? 0,
    timeZone: json['timeZone'] as String? ?? '',
  );

  final int from;
  final int to;
  final int discount;
  final String timeZone;
}
