class HardwareCheck {
  const HardwareCheck({
    required this.name,
    required this.required,
    required this.status,
    required this.detail,
  });

  factory HardwareCheck.fromJson(Map<String, Object?> json) => HardwareCheck(
    name: json['name']! as String,
    required: json['required'] as bool? ?? false,
    status: json['status'] as String? ?? off,
    detail: json['detail'] as String? ?? '',
  );

  static const ok = 'ok';
  static const off = 'off';
  static const failed = 'failed';

  final String name;
  final bool required;
  final String status;
  final String detail;
}
