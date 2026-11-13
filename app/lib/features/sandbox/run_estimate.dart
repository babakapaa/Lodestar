import 'estimate_option.dart';
import 'night_window.dart';
import 'sandbox_json.dart';

class RunEstimate {
  const RunEstimate({
    required this.maxHours,
    required this.starter,
    required this.available,
    required this.night,
    required this.options,
    required this.suggested,
  });

  factory RunEstimate.fromJson(Map<String, Object?> json) => RunEstimate(
    maxHours: SandboxJson.number(json['maxHours']),
    starter: json['starter'] as bool? ?? false,
    available: SandboxJson.number(json['available']),
    night: json['night'] == null
        ? null
        : NightWindow.fromJson(json['night']! as Map<String, Object?>),
    options: SandboxJson.list(json['options'], EstimateOption.fromJson),
    suggested: json['suggested'] as String?,
  );

  final double maxHours;
  final bool starter;
  final double available;
  final NightWindow? night;
  final List<EstimateOption> options;
  final String? suggested;
}
