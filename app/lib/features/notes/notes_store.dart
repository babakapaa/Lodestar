import 'dart:async';

import 'package:drift/drift.dart';
import 'package:lodestar_crdt/lodestar_crdt.dart';

import '../../core/offline/client_id.dart';
import 'note.dart';
import 'note_changes.dart';
import 'note_doc.dart';
import 'note_links.dart';
import 'note_meta.dart';
import 'note_role.dart';
import 'notes_database.dart';
import 'remote_note.dart';

class NotesStore {
  NotesStore(this._db, {this.flushDelay = const Duration(milliseconds: 400)});

  final NotesDatabase _db;
  final Duration flushDelay;
  final _docs = <String, Future<NoteDoc>>{};
  final _timers = <String, Timer>{};
  final _local = StreamController<String>.broadcast();
  Future<void> _queue = Future.value();

  Stream<String> get localChanges => _local.stream;

  Stream<List<Note>> watchAll() =>
      (_db.select(_db.noteRows)
            ..orderBy([(row) => OrderingTerm.desc(row.changedAt)]))
          .watch()
          .map((rows) => [for (final row in rows) Note.fromRow(row)]);

  Stream<Note?> watch(String id) =>
      (_db.select(_db.noteRows)..where((row) => row.id.equals(id)))
          .watchSingleOrNull()
          .map((row) => row == null ? null : Note.fromRow(row));

  Stream<bool> watchPending() =>
      _db.select(_db.pendingUpdates).watch().map((rows) => rows.isNotEmpty);

  Future<NoteRow?> row(String id) => (_db.select(
    _db.noteRows,
  )..where((row) => row.id.equals(id))).getSingleOrNull();

  Future<List<NoteRow>> rows() => _db.select(_db.noteRows).get();

  Future<String> create(NoteMeta meta) async {
    await CrdtText.ensureReady();
    final id = ClientId.generate();
    await _db
        .into(_db.noteRows)
        .insert(
          NoteRowsCompanion.insert(
            id: id,
            role: NoteRole.owner.name,
            path: meta.path,
            tags: meta.tags,
            links: meta.links,
            assignmentId: Value(meta.assignmentId),
            lessonId: Value(meta.lessonId),
            state: Uint8List(0),
            changedAt: DateTime.now().toUtc(),
          ),
        );
    _local.add(id);
    return id;
  }

  Future<void> setMeta(String id, NoteMeta meta) => _serial(() async {
    final current = await row(id);
    if (current == null) {
      return;
    }
    await _db
        .update(_db.noteRows)
        .replace(
          current.copyWith(
            path: meta.path,
            tags: meta.tags,
            links: meta.links,
            changedAt: DateTime.now().toUtc(),
            metaDirty: current.created,
            error: const Value(null),
          ),
        );
    _local.add(id);
  });

  Future<NoteDoc> open(String id) => _docs.putIfAbsent(id, () async {
    await CrdtText.ensureReady();
    final current = await row(id);
    return NoteDoc(id, current?.state ?? Uint8List(0));
  });

  void edit(NoteDoc doc, String text) {
    if (!doc.edit(text)) {
      return;
    }
    _timers[doc.id]?.cancel();
    _timers[doc.id] = Timer(flushDelay, () => unawaited(flush(doc.id)));
  }

  Future<void> flush(String id) => _serial(() => _flush(id));

  Future<void> flushAll() =>
      Future.wait([for (final id in _docs.keys) flush(id)]);

  Future<List<PendingUpdate>> pending() => (_db.select(
    _db.pendingUpdates,
  )..orderBy([(row) => OrderingTerm.asc(row.id)])).get();

  Future<void> sent(PendingUpdate update) => (_db.delete(
    _db.pendingUpdates,
  )..where((row) => row.id.equals(update.id))).go();

  Future<void> markCreated(String id) =>
      (_db.update(_db.noteRows)..where((row) => row.id.equals(id))).write(
        const NoteRowsCompanion(created: Value(true)),
      );

  Future<void> markMetaSent(String id, NoteMeta sent) => _serial(() async {
    final current = await row(id);
    if (current == null || current.path != sent.path) {
      return;
    }
    if (_sameList(current.tags, sent.tags) &&
        _sameList(current.links, sent.links)) {
      await (_db.update(_db.noteRows)..where((row) => row.id.equals(id))).write(
        const NoteRowsCompanion(metaDirty: Value(false)),
      );
    }
  });

  Future<void> markError(String id, String code) =>
      (_db.update(_db.noteRows)..where((row) => row.id.equals(id))).write(
        NoteRowsCompanion(error: Value(code)),
      );

  Future<void> receive(RemoteNote remote) => _serial(() async {
    final current = await row(remote.id);
    if (current == null) {
      await _db
          .into(_db.noteRows)
          .insert(
            NoteRowsCompanion.insert(
              id: remote.id,
              ownerId: Value(remote.ownerId),
              role: remote.role.name,
              path: remote.meta.path,
              tags: remote.meta.tags,
              links: remote.meta.links,
              assignmentId: Value(remote.meta.assignmentId),
              lessonId: Value(remote.meta.lessonId),
              state: Uint8List(0),
              changedAt: remote.changedAt,
              deletedAt: Value(remote.deletedAt),
              created: const Value(true),
            ),
          );
      return;
    }
    final keepMeta = current.metaDirty;
    await _db
        .update(_db.noteRows)
        .replace(
          current.copyWith(
            ownerId: remote.ownerId,
            role: remote.role.name,
            path: keepMeta ? current.path : remote.meta.path,
            tags: keepMeta ? current.tags : remote.meta.tags,
            links: keepMeta ? current.links : remote.meta.links,
            assignmentId: Value(remote.meta.assignmentId),
            lessonId: Value(remote.meta.lessonId),
            changedAt: remote.changedAt.isAfter(current.changedAt)
                ? remote.changedAt
                : current.changedAt,
            deletedAt: Value(remote.deletedAt),
            created: true,
          ),
        );
  });

  Future<void> applyChanges(String id, NoteChanges changes) =>
      _serial(() async {
        await _flush(id);
        final doc = await open(id);
        doc.applyRemote(changes.all);
        final current = await row(id);
        if (current == null) {
          return;
        }
        await _db
            .update(_db.noteRows)
            .replace(
              current.copyWith(
                state: doc.state,
                body: doc.text,
                seq: changes.seq,
                unsnapshotted: changes.snapshot == null
                    ? current.unsnapshotted + changes.updates.length
                    : changes.updates.length,
              ),
            );
      });

  Future<void> markSnapshotted(String id) =>
      (_db.update(_db.noteRows)..where((row) => row.id.equals(id))).write(
        const NoteRowsCompanion(unsnapshotted: Value(0)),
      );

  Future<void> forget(String id) => _serial(() async {
    _timers.remove(id)?.cancel();
    (await _docs.remove(id))?.dispose();
    await _db.transaction(() async {
      await (_db.delete(
        _db.pendingUpdates,
      )..where((row) => row.noteId.equals(id))).go();
      await (_db.delete(_db.noteRows)..where((row) => row.id.equals(id))).go();
    });
  });

  Future<void> clear() => _serial(() async {
    for (final timer in _timers.values) {
      timer.cancel();
    }
    _timers.clear();
    for (final doc in _docs.values) {
      (await doc).dispose();
    }
    _docs.clear();
    await _db.clear();
  });

  Future<void> dispose() async {
    for (final timer in _timers.values) {
      timer.cancel();
    }
    _timers.clear();
    await _local.close();
    for (final doc in _docs.values) {
      (await doc).dispose();
    }
    _docs.clear();
  }

  Future<void> _flush(String id) async {
    _timers.remove(id)?.cancel();
    final opened = _docs[id];
    if (opened == null) {
      return;
    }
    final doc = await opened;
    final update = doc.takeLocal();
    if (update == null) {
      return;
    }
    final current = await row(id);
    if (current == null) {
      return;
    }
    final body = doc.text;
    final links = NoteLinks.parse(body);
    final linksChanged = !_sameList(links, current.links);
    await _db.transaction(() async {
      await _db
          .into(_db.pendingUpdates)
          .insert(
            PendingUpdatesCompanion.insert(
              noteId: id,
              key: ClientId.generate(),
              data: update,
            ),
          );
      await _db
          .update(_db.noteRows)
          .replace(
            current.copyWith(
              state: doc.state,
              body: body,
              links: links,
              changedAt: DateTime.now().toUtc(),
              metaDirty: current.metaDirty || (linksChanged && current.created),
            ),
          );
    });
    _local.add(id);
  }

  Future<T> _serial<T>(Future<T> Function() action) {
    final result = _queue.then((_) => action());
    _queue = result.then<void>((_) {}, onError: (_) {});
    return result;
  }

  static bool _sameList(List<String> a, List<String> b) {
    if (a.length != b.length) {
      return false;
    }
    for (var i = 0; i < a.length; i++) {
      if (a[i] != b[i]) {
        return false;
      }
    }
    return true;
  }
}
