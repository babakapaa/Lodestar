import 'package:intl/intl.dart';

import '../../l10n/generated/app_localizations.dart';
import 'run_kind.dart';
import 'run_status.dart';
import 'sandbox_run.dart';

class SandboxLabels {
  const SandboxLabels(this._l10n, this._locale);

  final AppLocalizations _l10n;
  final String _locale;

  String hours(double value) => _l10n.sandboxHours(
    NumberFormat.decimalPatternDigits(
      locale: _locale,
      decimalDigits: value == value.roundToDouble() ? 0 : 2,
    ).format(value),
  );

  String kind(RunKind kind) => switch (kind) {
    RunKind.jupyter => _l10n.sandboxKindJupyter,
    RunKind.vscode => _l10n.sandboxKindVscode,
    RunKind.job => _l10n.sandboxKindJob,
    RunKind.demo => _l10n.sandboxKindDemo,
  };

  String status(RunStatus status) => switch (status) {
    RunStatus.queued => _l10n.sandboxStatusQueued,
    RunStatus.starting => _l10n.sandboxStatusStarting,
    RunStatus.running => _l10n.sandboxStatusRunning,
    RunStatus.held => _l10n.sandboxStatusHeld,
    RunStatus.stopping => _l10n.sandboxStatusStopping,
    RunStatus.finished => _l10n.sandboxStatusFinished,
    RunStatus.failed => _l10n.sandboxStatusFailed,
    RunStatus.stopped => _l10n.sandboxStatusStopped,
  };

  String? reason(String reason) => switch (reason) {
    'time_limit' => _l10n.sandboxReasonTimeLimit,
    'node_lost' => _l10n.sandboxReasonNodeLost,
    'vanished' => _l10n.sandboxReasonVanished,
    'idle' => _l10n.sandboxReasonIdle,
    'watchdog' => _l10n.sandboxReasonWatchdog,
    _ => null,
  };

  String source(String source, String title) {
    final name = switch (source) {
      'sandbox.base' => _l10n.sandboxSourceBase,
      'sandbox.topup' => _l10n.sandboxSourceTopup,
      'sandbox.transfer' => _l10n.sandboxSourceTransfer,
      'sandbox.run' => _l10n.sandboxSourceRun,
      'league.division' => _l10n.sandboxSourceDivision,
      'league.award' => _l10n.sandboxSourceAward,
      _ => source,
    };
    return title.isEmpty ? name : '$name: $title';
  }

  String summary(SandboxRun run) => [
    status(run.status),
    if (run.position case final position?) _l10n.sandboxPosition(position),
    run.className,
    if (run.free)
      _l10n.sandboxFree
    else if (run.charged case final charged?)
      hours(charged),
  ].join(', ');
}
