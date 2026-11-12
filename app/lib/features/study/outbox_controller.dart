import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/api_exception.dart';
import '../../core/auth/auth_providers.dart';
import '../../core/auth/auth_state.dart';
import '../../core/live/live_providers.dart';
import 'study_providers.dart';
import 'submission_draft.dart';

class OutboxController extends Notifier<List<SubmissionDraft>> {
  static const _firstDelay = Duration(seconds: 5);
  static const _maxDelay = Duration(seconds: 60);

  Timer? _timer;
  Duration _delay = _firstDelay;
  bool _sending = false;

  @override
  List<SubmissionDraft> build() {
    ref.onDispose(() => _timer?.cancel());
    ref.listen(liveEventsProvider, (_, _) => flush());
    ref.listen(authProvider, (_, next) {
      if (next is Disconnected) {
        _save(const []);
      }
    });
    final drafts = ref.read(submissionOutboxProvider).load();
    if (drafts.any((d) => !d.failed)) {
      Future.microtask(flush);
    }
    return drafts;
  }

  Future<void> add(SubmissionDraft draft) async {
    await _save([...state, draft]);
    await flush();
  }

  Future<void> discard(String clientId) => _save([
    for (final d in state)
      if (d.clientId != clientId) d,
  ]);

  Future<void> retry(String clientId) async {
    await _save([
      for (final d in state) d.clientId == clientId ? d.withError(null) : d,
    ]);
    await flush();
  }

  Future<void> flush() async {
    if (_sending) {
      return;
    }
    _sending = true;
    _timer?.cancel();
    try {
      for (final draft in [...state]) {
        if (draft.failed) {
          continue;
        }
        try {
          await ref.read(studyRepositoryProvider).submit(draft);
          await _save([
            for (final d in state)
              if (d.clientId != draft.clientId) d,
          ]);
          ref.invalidate(assignmentProvider(draft.assignmentId));
          ref.invalidate(assignmentsProvider);
        } on ApiException catch (error) {
          if (error.code == ApiException.unreachable ||
              error.statusCode == 401 ||
              error.statusCode >= 500) {
            _schedule();
            return;
          }
          await _save([
            for (final d in state)
              d.clientId == draft.clientId ? d.withError(error.code) : d,
          ]);
        }
      }
      _delay = _firstDelay;
    } finally {
      _sending = false;
    }
  }

  void _schedule() {
    _timer = Timer(_delay, flush);
    final next = _delay * 2;
    _delay = next > _maxDelay ? _maxDelay : next;
  }

  Future<void> _save(List<SubmissionDraft> drafts) async {
    state = drafts;
    await ref.read(submissionOutboxProvider).save(drafts);
  }
}
