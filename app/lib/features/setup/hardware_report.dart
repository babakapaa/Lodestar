import 'hardware_check.dart';

class HardwareReport {
  const HardwareReport({required this.checkedAt, required this.checks});

  factory HardwareReport.fromJson(Map<String, Object?> json) => HardwareReport(
    checkedAt: DateTime.tryParse(json['checkedAt'] as String? ?? '')?.toLocal(),
    checks: [
      for (final item in json['checks'] as List<Object?>? ?? const [])
        HardwareCheck.fromJson(item! as Map<String, Object?>),
    ],
  );

  final DateTime? checkedAt;
  final List<HardwareCheck> checks;
}
