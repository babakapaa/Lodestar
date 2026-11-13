import '../../l10n/generated/app_localizations.dart';
import '../study/points_text.dart';
import '../study/study_labels.dart';
import 'dispute_status.dart';
import 'entry_kind.dart';
import 'league_display.dart';
import 'league_entry.dart';
import 'league_transition.dart';
import 'standing.dart';
import 'transition_status.dart';

class LeagueLabels {
  const LeagueLabels(this._l10n, this._study);

  final AppLocalizations _l10n;
  final StudyLabels _study;

  String display(LeagueDisplay display) => switch (display) {
    LeagueDisplay.name => _l10n.leagueDisplayName,
    LeagueDisplay.hidden => _l10n.leagueDisplayHidden,
    LeagueDisplay.nickname => _l10n.leagueDisplayNickname,
  };

  String toNext(Standing standing, String next) => standing.toNext == null
      ? _l10n.leagueTopDivision
      : _l10n.leagueToNext(standing.toNext!, next);

  String growth(int growth) =>
      growth > 0 ? _l10n.leagueGrowth(growth) : _l10n.leagueNoGrowth;

  String entry(LeagueEntry entry) => switch (entry.kind) {
    EntryKind.work when entry.score != null && entry.max != null =>
      _l10n.leagueEntryWork(
        PointsText.of(entry.score!),
        PointsText.of(entry.max!),
        entry.weight ?? 0,
      ),
    EntryKind.work => _l10n.leagueEntryPending,
    EntryKind.adjustment => _l10n.leagueEntryAdjustment(entry.reason),
    EntryKind.task => _l10n.leagueEntryTask,
  };

  String points(int points) =>
      points > 0 ? _l10n.leaguePlus(points) : _l10n.leaguePoints(points);

  String transition(LeagueTransition t, String from, String to) =>
      t.up ? _l10n.leagueMoveUp(from, to) : _l10n.leagueMoveDown(from, to);

  String transitionStatus(TransitionStatus status) => switch (status) {
    TransitionStatus.pending => _l10n.leagueMovePending,
    TransitionStatus.applied => _l10n.leagueMoveApplied,
    TransitionStatus.rejected => _l10n.leagueMoveRejected,
    TransitionStatus.reverted => _l10n.leagueMoveReverted,
  };

  String disputeStatus(DisputeStatus status) => switch (status) {
    DisputeStatus.open => _l10n.leagueDisputeOpen,
    DisputeStatus.accepted => _l10n.leagueDisputeAccepted,
    DisputeStatus.rejected => _l10n.leagueDisputeRejected,
  };

  String time(DateTime at) => _study.time(at);

  String due(DateTime? at) => _study.due(at);
}
