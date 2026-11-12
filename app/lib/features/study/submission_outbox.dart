import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

import 'submission_draft.dart';

class SubmissionOutbox {
  const SubmissionOutbox(this._prefs);

  static const _key = 'outbox.submissions';

  final SharedPreferences _prefs;

  List<SubmissionDraft> load() {
    final raw = _prefs.getString(_key);
    if (raw == null) {
      return const [];
    }
    try {
      return [
        for (final item in jsonDecode(raw) as List<Object?>)
          SubmissionDraft.fromJson(item! as Map<String, Object?>),
      ];
    } on Object {
      return const [];
    }
  }

  Future<void> save(List<SubmissionDraft> drafts) => _prefs.setString(
    _key,
    jsonEncode([for (final draft in drafts) draft.toJson()]),
  );
}
