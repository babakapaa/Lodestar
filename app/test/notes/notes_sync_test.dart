import 'package:drift/drift.dart' show driftRuntimeOptions;
import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/api/api_exception.dart';
import 'package:lodestar/features/notes/note_meta.dart';

import '../support/fake_server.dart';
import 'notes_device.dart';

void main() {
  late FakeServer server;
  late NotesDevice aigerim;
  late NotesDevice erlan;

  setUpAll(() => driftRuntimeOptions.dontWarnAboutMultipleDatabases = true);

  setUp(() {
    server = FakeServer();
    aigerim = NotesDevice(server);
    erlan = NotesDevice(server);
  });

  tearDown(() async {
    await aigerim.close();
    await erlan.close();
  });

  test('two people edit one note offline and nothing is lost', () async {
    final id = await aigerim.create('Биология/Клетка');
    await aigerim.type(id, 'Клетка состоит из ядра.');
    await aigerim.sync.run();
    await erlan.sync.run();
    expect(await erlan.text(id), 'Клетка состоит из ядра.');

    server.offline = true;
    await aigerim.type(id, 'Айгерим: Клетка состоит из ядра.');
    await erlan.type(id, 'Клетка состоит из ядра. И мембраны 🧬.');
    await expectLater(
      aigerim.sync.run(),
      throwsA(
        isA<ApiException>().having(
          (e) => e.code,
          'code',
          ApiException.unreachable,
        ),
      ),
    );

    server.offline = false;
    await aigerim.sync.run();
    await erlan.sync.run();
    await aigerim.sync.run();

    const merged = 'Айгерим: Клетка состоит из ядра. И мембраны 🧬.';
    expect(await aigerim.text(id), merged);
    expect(await erlan.text(id), merged);
    expect((await erlan.notes()).single.body, merged);
  });

  test('a note created offline reaches the server with its meta', () async {
    server.offline = true;
    final id = await aigerim.create('Черновик');
    await aigerim.type(id, 'См. [[Клетка]] и [[Биология/Ткани|ткани]]');
    server.offline = false;

    await aigerim.sync.run();

    final meta = server.notes.notes[id]!['meta']! as Map<String, Object?>;
    expect(meta['path'], 'Черновик');
    expect(meta['links'], ['Клетка', 'Биология/Ткани']);
    expect(server.notes.updates[id], hasLength(1));
    expect((await aigerim.notes()).single.error, isNull);
  });

  test('renaming keeps local meta until it is sent', () async {
    final id = await aigerim.create('Клетка');
    await aigerim.sync.run();
    await aigerim.store.setMeta(
      id,
      const NoteMeta(path: 'Биология/Клетка', tags: ['экзамен']),
    );

    await aigerim.sync.run();

    final note = (await aigerim.notes()).single;
    expect(note.meta.path, 'Биология/Клетка');
    expect(note.meta.tags, ['экзамен']);
    expect(
      (server.notes.notes[id]!['meta']! as Map<String, Object?>)['path'],
      'Биология/Клетка',
    );
  });

  test(
    'many updates are folded into a snapshot that others can read',
    () async {
      final id = await aigerim.create('Лекция');
      await aigerim.sync.run();
      for (final text in ['a', 'ab', 'abc', 'abcd']) {
        await erlan.sync.run();
        await aigerim.type(id, text);
        await aigerim.sync.run();
      }
      await erlan.sync.run();
      expect(server.notes.snapshots[id], isNotNull);

      final newcomer = NotesDevice(server);
      addTearDown(newcomer.close);
      await newcomer.sync.run();
      expect(await newcomer.text(id), 'abcd');
    },
  );

  test('a note the server no longer shares disappears locally', () async {
    final id = await aigerim.create('Общая');
    await aigerim.sync.run();
    await erlan.sync.run();
    expect(await erlan.notes(), hasLength(1));

    server.notes.notes.remove(id);
    await erlan.sync.run();

    expect(await erlan.notes(), isEmpty);
  });
}
