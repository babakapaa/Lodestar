import '../../l10n/generated/app_localizations.dart';
import 'sandbox_labels.dart';
import 'sandbox_run.dart';

class RunQuestion {
  const RunQuestion(this._l10n, this._labels);

  static const logLines = 40;

  final AppLocalizations _l10n;
  final SandboxLabels _labels;

  String of(SandboxRun run) {
    final lines = run.log.trimRight().split('\n');
    final tail = lines.length > logLines
        ? lines.sublist(lines.length - logLines)
        : lines;
    return [
      _l10n.sandboxWhyPrompt,
      if (run.command.isNotEmpty) '${_l10n.sandboxCommand}: ${run.command}',
      if (run.exitCode case final code?) '${_l10n.sandboxFactExit}: $code',
      ?_labels.reason(run.reason),
      if (run.log.trim().isNotEmpty) '```\n${tail.join('\n')}\n```',
    ].join('\n');
  }
}
