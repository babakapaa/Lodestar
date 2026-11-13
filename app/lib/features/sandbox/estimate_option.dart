import 'sandbox_json.dart';

class EstimateOption {
  const EstimateOption({
    required this.classId,
    required this.className,
    required this.memoryGb,
    required this.rate,
    required this.cost,
    required this.nightCost,
    required this.freeGpus,
    required this.queued,
    required this.enough,
  });

  factory EstimateOption.fromJson(Map<String, Object?> json) => EstimateOption(
    classId: json['classId']! as String,
    className: json['className'] as String? ?? '',
    memoryGb: json['memoryGb'] as int? ?? 0,
    rate: SandboxJson.number(json['rate']),
    cost: SandboxJson.number(json['cost']),
    nightCost: SandboxJson.number(json['nightCost']),
    freeGpus: json['freeGpus'] as int? ?? 0,
    queued: json['queued'] as int? ?? 0,
    enough: json['enough'] as bool? ?? false,
  );

  final String classId;
  final String className;
  final int memoryGb;
  final double rate;
  final double cost;
  final double nightCost;
  final int freeGpus;
  final int queued;
  final bool enough;

  bool get startsNow => freeGpus > queued;
}
