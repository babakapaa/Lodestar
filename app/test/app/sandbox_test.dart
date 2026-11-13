import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../support/fake_server.dart';
import 'admin_navigation.dart';
import 'app_harness.dart';

Future<void> go(WidgetTester tester, String location) async {
  GoRouter.of(tester.element(find.byType(Scaffold).first)).go(location);
  await tester.pumpAndSettle();
}

FakeServer sandboxServer() {
  final server = FakeServer();
  server.modules.add({
    'name': 'sandbox',
    'version': '0.1.0',
    'optional': true,
    'enabled': true,
  });
  return server;
}

void main() {
  testWidgets('a student sees the balance and starts a notebook', (
    tester,
  ) async {
    final server = sandboxServer();
    final app = await pumpApp(
      tester,
      size: const Size(390, 2400),
      server: server,
    );
    await tester.tap(find.text('Учёба').first);
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('my-sandbox')));
    await tester.pumpAndSettle();

    expect(find.text('26 ч'), findsOneWidget);
    expect(
      find.text('При нынешнем темпе к концу семестра останется 30 ч'),
      findsOneWidget,
    );
    expect(find.text('Нейроны'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('give-team-1')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('sandbox-hours')), '1,5');
    await tester.tap(find.byKey(const ValueKey('sandbox-hours-done')));
    await tester.pumpAndSettle();
    expect(server.sandbox.transfers.single, {'teamId': 'team-1', 'hours': 1.5});

    await tester.tap(find.byKey(const ValueKey('sandbox-export')));
    await tester.pumpAndSettle();
    await tester.tap(find.text('CSV'));
    await tester.pumpAndSettle();
    expect(app.saver.saved.keys.single, endsWith('.csv'));

    await tester.tap(find.byKey(const ValueKey('sandbox-new')));
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('run-start')), findsOneWidget);
    await tester.enterText(find.byKey(const ValueKey('run-title')), 'Ноутбук');
    await tester.tap(find.byKey(const ValueKey('run-template')));
    await tester.pumpAndSettle();
    await tester.tap(find.text('PyTorch').last);
    await tester.pumpAndSettle();

    expect(find.text('Можно потратить 26 ч'), findsOneWidget);
    expect(find.text('С 23:00 до 7:00 скидка 50%.'), findsOneWidget);
    expect(find.textContaining('до 8 ч, ночью 4 ч'), findsOneWidget);
    expect(find.textContaining('не хватает часов'), findsOneWidget);

    await tester.ensureVisible(find.byKey(const ValueKey('run-start')));
    await tester.tap(find.byKey(const ValueKey('run-start')));
    await tester.pumpAndSettle();
    expect(server.sandbox.started.single, {
      'kind': 'jupyter',
      'title': 'Ноутбук',
      'templateId': 'tp1',
      'classId': 'c-t4',
    });
    expect(find.text('В очереди, 1-й в очереди, T4'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('run-stop')));
    await tester.pumpAndSettle();
    expect(server.sandbox.stopped, ['r1']);
  });

  testWidgets('a job for course work goes with the assignment', (tester) async {
    final server = sandboxServer();
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await go(
      tester,
      '/sandbox/new?workKind=assignment&workId=a1&workTitle=МРТ',
    );
    expect(find.text('Для задания «МРТ»'), findsOneWidget);
    await tester.tap(find.text('Задача'));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey('run-template')));
    await tester.pumpAndSettle();
    await tester.tap(find.text('PyTorch').last);
    await tester.pumpAndSettle();
    expect(find.text('Добавьте команду'), findsOneWidget);
    await tester.enterText(
      find.byKey(const ValueKey('run-command')),
      'python train.py',
    );
    await tester.enterText(find.byKey(const ValueKey('run-hours')), '2');
    await tester.pumpAndSettle();
    await tester.ensureVisible(find.byKey(const ValueKey('run-start')));
    await tester.tap(find.byKey(const ValueKey('run-start')));
    await tester.pumpAndSettle();
    expect(server.sandbox.started.single, {
      'kind': 'job',
      'title': 'МРТ',
      'templateId': 'tp1',
      'classId': 'c-t4',
      'command': 'python train.py',
      'maxHours': 2.0,
      'workKind': 'assignment',
      'workId': 'a1',
    });
  });

  testWidgets('a demo shows its link and QR code', (tester) async {
    final server = sandboxServer();
    server.sandbox.runs.add(server.sandbox.run('d1', 'demo', 'running'));
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await go(tester, '/sandbox/runs/d1');

    expect(find.byKey(const ValueKey('run-open')), findsOneWidget);
    expect(find.byKey(const ValueKey('run-qr')), findsOneWidget);
  });

  testWidgets('a failed run explains itself and runs again', (tester) async {
    final server = sandboxServer();
    server.sandbox.runs.add(
      server.sandbox.run('f1', 'job', 'failed')
        ..['refunded'] = true
        ..['charged'] = 0
        ..['exitCode'] = 1
        ..['log'] = 'ModuleNotFoundError: torchvision',
    );
    server.sandbox.runs.add(
      server.sandbox.run('h1', 'job', 'held')..['holdNote'] = 'xmrig',
    );
    server.sandbox.runs.add(
      server.sandbox.run('w1', 'job', 'finished')
        ..['charged'] = 3
        ..['gpuLoad'] = 12,
    );
    await pumpApp(tester, size: const Size(390, 2400), server: server);

    await go(tester, '/sandbox/runs/h1');
    expect(find.textContaining('найден процесс xmrig'), findsOneWidget);

    await go(tester, '/sandbox/runs/w1');
    expect(find.textContaining('загружена в среднем на 12%'), findsOneWidget);

    await go(tester, '/sandbox/runs/f1');
    expect(
      find.text('Запуск упал в первые 5 минут, часы вернулись.'),
      findsOneWidget,
    );
    expect(find.text('ModuleNotFoundError: torchvision'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('run-ask')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('sandbox-hours')), '5');
    await tester.enterText(
      find.byKey(const ValueKey('sandbox-reason')),
      'Обучение дольше',
    );
    await tester.tap(find.byKey(const ValueKey('sandbox-hours-done')));
    await tester.pumpAndSettle();
    expect(server.sandbox.topups.single['hours'], 5);
    expect(server.sandbox.topups.single['workId'], 'a1');

    await tester.tap(find.byKey(const ValueKey('run-rerun')));
    await tester.pumpAndSettle();
    expect(server.sandbox.runs.first['title'], 'Запуск f1');
    expect(find.text('В очереди, 1-й в очереди, T4'), findsOneWidget);
  });

  testWidgets('a teacher adds requested hours', (tester) async {
    final server = sandboxServer();
    server.sandbox.topups.add({
      'id': 'tu9',
      'userId': 'u2',
      'displayName': 'Ерлан',
      'workKind': 'study.assignment',
      'workId': 'a1',
      'title': 'Нейросеть для МРТ',
      'hours': 6,
      'reason': 'Большой датасет',
      'status': 'pending',
      'mine': false,
      'createdAt': '2026-11-13T08:00:00Z',
    });
    await pumpApp(tester, size: const Size(390, 2400), server: server);
    await go(tester, '/sandbox/topups');

    expect(find.text('Ерлан, Нейросеть для МРТ'), findsOneWidget);
    await tester.tap(find.byKey(const ValueKey('approve-tu9')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('sandbox-hours')), '4');
    await tester.tap(find.byKey(const ValueKey('sandbox-hours-done')));
    await tester.pumpAndSettle();
    expect(server.sandbox.decisions.single, {
      'topup': 'tu9',
      'approve': true,
      'hours': 4,
    });
  });

  testWidgets('an administrator connects a GPU server and decides a hold', (
    tester,
  ) async {
    final server = sandboxServer()..permissions = ['sandbox.manage'];
    server.sandbox.runs.add(
      server.sandbox.run('h1', 'job', 'held')..['holdNote'] = 'xmrig',
    );
    await pumpApp(tester, size: const Size(390, 2600), server: server);
    await openAdmin(tester, 'Семестры, видеокарты, окружения и серверы');

    expect(
      find.text('Найден процесс xmrig. Время на паузе не списывается.'),
      findsOneWidget,
    );
    await tester.tap(find.byKey(const ValueKey('held-release-h1')));
    await tester.pumpAndSettle();
    expect(server.sandbox.decisions.single, {'run': 'h1', 'release': true});

    await tester.ensureVisible(find.byKey(const ValueKey('node-new')));
    await tester.tap(find.byKey(const ValueKey('node-new')));
    await tester.pumpAndSettle();
    await tester.enterText(find.byKey(const ValueKey('node-name')), 'gpu-01');
    await tester.enterText(find.byKey(const ValueKey('node-gpus')), '4');
    await tester.tap(find.byKey(const ValueKey('node-save')));
    await tester.pumpAndSettle();
    expect(find.text('lsn_t1_secret'), findsOneWidget);
    await tester.tap(find.byKey(const ValueKey('node-token-done')));
    await tester.pumpAndSettle();
    expect(server.sandbox.nodes.single['gpus'], 4);
    expect(find.text('gpu-01'), findsOneWidget);
  });
}
