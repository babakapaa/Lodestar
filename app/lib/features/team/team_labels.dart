import '../../l10n/generated/app_localizations.dart';
import '../study/study_labels.dart';
import 'history_entry.dart';
import 'history_event.dart';
import 'team_scope.dart';
import 'vote_decision.dart';

class TeamLabels {
  const TeamLabels(this._l10n, this._study);

  final AppLocalizations _l10n;
  final StudyLabels _study;

  String size(TeamScope scope) => scope.minSize == scope.maxSize
      ? _l10n.teamSizeExact(scope.minSize)
      : _l10n.teamSizeRange(scope.minSize, scope.maxSize);

  String rules(TeamScope scope) => [
    size(scope),
    if (scope.selfFormed) _l10n.teamSelfFormed else _l10n.teamByTeacher,
    if (scope.freeLeaveUntil case final until?)
      _l10n.teamFreeLeaveUntil(_study.time(until)),
    _study.due(scope.dueAt),
  ].join(', ');

  String event(HistoryEntry entry) {
    final who = entry.person?.label ?? _l10n.teamSomeone;
    final actor = entry.actor?.label ?? _l10n.teamSomeone;
    return switch (entry.event) {
      HistoryEvent.created => _l10n.teamEventCreated(actor),
      HistoryEvent.joined => _l10n.teamEventJoined(who),
      HistoryEvent.added => _l10n.teamEventAdded(who, actor),
      HistoryEvent.left => _l10n.teamEventLeft(who),
      HistoryEvent.removed => _l10n.teamEventRemoved(who, actor),
      HistoryEvent.separated => _l10n.teamEventSeparated(who, actor),
      HistoryEvent.movedIn => _l10n.teamEventMovedIn(who, actor),
      HistoryEvent.movedOut => _l10n.teamEventMovedOut(who, actor),
    };
  }

  String decision(VoteDecision decision) => switch (decision) {
    VoteDecision.waiting => _l10n.voteWaiting,
    VoteDecision.agreed => _l10n.voteAgreed,
    VoteDecision.objected => _l10n.voteObjected,
    VoteDecision.silent => _l10n.voteSilent,
  };

  String time(DateTime at) => _study.time(at);
}
