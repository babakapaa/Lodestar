import 'package:flutter_test/flutter_test.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'app_harness.dart';

Map<String, Object?> lesson(String date, String kind, String status) => {
  'seriesId': 'sr-$kind',
  'courseId': 'c1',
  'courseTitle': 'Биология',
  'unitId': null,
  'kind': kind,
  'date': date,
  'startsAt': '${date}T04:00:00Z',
  'endsAt': '${date}T05:30:00Z',
  'room': 'А-204',
  'teacherId': null,
  'status': status,
  'note': status == 'cancelled' ? 'Преподаватель болеет' : '',
};

void main() {
  testWidgets('today shows classes and deadlines and opens the week', (
    tester,
  ) async {
    final server = FakeServer();
    server.study.lessons.addAll([
      lesson('2026-11-12', 'lecture', 'planned'),
      lesson('2026-11-13', 'practice', 'cancelled'),
    ]);
    server.study.deadlines.add({
      'assignmentId': 'a1',
      'courseId': 'c1',
      'courseTitle': 'Биология',
      'title': 'Эссе о клетке',
      'dueAt': '2026-11-12T13:00:00Z',
      'submittedAt': null,
    });
    await pumpApp(tester, size: const Size(390, 1400), server: server);

    expect(find.text('Занятия'), findsOneWidget);
    expect(find.text('Сроки сдачи'), findsOneWidget);
    expect(find.text('Эссе о клетке'), findsOneWidget);
    expect(find.textContaining('Лекция, А-204'), findsWidgets);

    await tester.tap(find.text('Расписание на неделю'));
    await tester.pumpAndSettle();
    expect(find.text('Расписание'), findsOneWidget);
    expect(
      find.textContaining('Практика, А-204, Отменено, Преподаватель болеет'),
      findsOneWidget,
    );
    final request = server.requests.lastWhere(
      (r) => r.url.path == '/v1/study/schedule',
    );
    expect(request.url.queryParameters['from'], '2026-11-09');
    expect(request.url.queryParameters['to'], '2026-11-15');

    await tester.tap(find.byTooltip('Следующая неделя'));
    await tester.pumpAndSettle();
    expect(server.requests.last.url.queryParameters['from'], '2026-11-16');
  });

  testWidgets('a notification opens the assignment it is about', (
    tester,
  ) async {
    final server = FakeServer();
    server.notifications.first['link'] = '/study/assignments/a1';
    await pumpApp(tester, size: const Size(390, 1400), server: server);

    await tester.tap(find.byTooltip('Уведомления, непрочитанных: 1'));
    await tester.pumpAndSettle();
    await tester.tap(find.textContaining('Скоро срок'));
    await tester.pumpAndSettle();

    expect(find.text('Эссе о клетке'), findsOneWidget);
    expect(server.notifications.first['readAt'], isNotNull);
  });
}
