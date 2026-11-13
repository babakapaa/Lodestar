import '../../l10n/generated/app_localizations.dart';
import 'unit_metric.dart';

class MetricLabels {
  const MetricLabels(this.l10n);

  final AppLocalizations l10n;

  String name(String key) => switch (key) {
    'study.submitted' => l10n.metricSubmitted,
    'study.on_time' => l10n.metricOnTime,
    'attendance.rate' => l10n.metricAttendance,
    'navigator.users' => l10n.metricNavigator,
    'sandbox.hours' => l10n.metricGpu,
    _ => key,
  };

  bool isShare(String key) => key != 'sandbox.hours';

  String value(UnitMetric metric) => isShare(metric.key)
      ? '${(metric.value * 100).round()}%'
      : l10n.sandboxHours(_number(metric.value));

  String _number(double value) => value == value.roundToDouble()
      ? value.toStringAsFixed(0)
      : value.toStringAsFixed(1);
}
