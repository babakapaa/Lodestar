import 'sandbox_json.dart';

class QuotaLine {
  const QuotaLine({
    required this.source,
    required this.title,
    required this.hours,
  });

  factory QuotaLine.fromJson(Map<String, Object?> json) => QuotaLine(
    source: json['source'] as String? ?? '',
    title: json['title'] as String? ?? '',
    hours: SandboxJson.number(json['hours']),
  );

  final String source;
  final String title;
  final double hours;
}
