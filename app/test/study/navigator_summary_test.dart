import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/study/assignment.dart';
import 'package:lodestar/features/study/submission.dart';

void main() {
  test('a submission carries the navigator summary', () {
    final submission = Submission.fromJson({
      'id': 'w1',
      'body': 'Ответ',
      'submittedAt': '2026-11-19T10:00:00Z',
      'files': <Object?>[],
      'searchMap': null,
      'navigatorSummary': {
        'replies': 3,
        'steps': {'question': 1, 'direction': 1, 'source': 1},
        'sources': ['Иванов И. И.'],
        'articles': <Object?>[],
      },
    });
    final help = submission.help!;
    expect(help.replies, 3);
    expect(help.steps.source, 1);
    expect(help.steps.example, 0);
    expect(help.sources, ['Иванов И. И.']);
  });

  test('the summary is off unless the lecturer turns it on', () {
    final assignment = Assignment.fromJson({
      'id': 'a1',
      'courseId': 'c1',
      'title': 'Мембрана',
    });
    expect(assignment.navigatorSummary, isFalse);
  });
}
