import '../../core/api/api_exception.dart';
import 'note_meta.dart';
import 'notes_api.dart';
import 'notes_database.dart';
import 'notes_store.dart';
import 'remote_note.dart';

class NotesSync {
  const NotesSync(this._api, this._store, {this.compactAfter = 100});

  final NotesApi _api;
  final NotesStore _store;
  final int compactAfter;

  Future<void> run() async {
    await _store.flushAll();
    await _sendCreated();
    await _sendMeta();
    await _sendUpdates();
    final remote = await _api.list();
    for (final note in remote) {
      await _store.receive(note);
    }
    for (final note in remote) {
      await _pull(note);
    }
    await _forgetMissing({for (final note in remote) note.id});
  }

  Future<void> saveVersion(String id) async {
    await run();
    final row = await _store.row(id);
    if (row == null) {
      return;
    }
    final doc = await _store.open(id);
    await _api.compact(id, row.seq, doc.state, doc.text);
    await _store.markSnapshotted(id);
  }

  Future<void> _sendCreated() async {
    for (final row in await _store.rows()) {
      if (row.created) {
        continue;
      }
      await _guard(row.id, () async {
        await _api.create(row.id, _meta(row));
        await _store.markCreated(row.id);
      });
    }
  }

  Future<void> _sendMeta() async {
    for (final row in await _store.rows()) {
      if (!row.created || !row.metaDirty) {
        continue;
      }
      final meta = _meta(row);
      await _guard(row.id, () async {
        await _api.setMeta(row.id, meta);
        await _store.markMetaSent(row.id, meta);
      });
    }
  }

  Future<void> _sendUpdates() async {
    final created = {
      for (final row in await _store.rows())
        if (row.created) row.id,
    };
    for (final update in await _store.pending()) {
      if (!created.contains(update.noteId)) {
        continue;
      }
      await _guard(update.noteId, () async {
        await _api.push(update.noteId, update.key, update.data);
        await _store.sent(update);
      });
    }
  }

  Future<void> _pull(RemoteNote note) async {
    final row = await _store.row(note.id);
    if (row == null || note.seq <= row.seq) {
      return;
    }
    await _guard(note.id, () async {
      await _store.applyChanges(note.id, await _api.changes(note.id, row.seq));
    });
    final pulled = await _store.row(note.id);
    if (pulled == null ||
        !note.role.canEdit ||
        pulled.unsnapshotted < compactAfter) {
      return;
    }
    await _guard(note.id, () async {
      final doc = await _store.open(note.id);
      await _api.compact(note.id, pulled.seq, doc.state, doc.text);
      await _store.markSnapshotted(note.id);
    });
  }

  Future<void> _forgetMissing(Set<String> remote) async {
    for (final row in await _store.rows()) {
      if (row.created && !remote.contains(row.id)) {
        await _store.forget(row.id);
      }
    }
  }

  Future<void> _guard(String id, Future<void> Function() action) async {
    try {
      await action();
    } on ApiException catch (error) {
      if (offline(error)) {
        rethrow;
      }
      await _store.markError(id, error.code);
    }
  }

  static bool offline(ApiException error) =>
      error.code == ApiException.unreachable ||
      error.statusCode == 401 ||
      error.statusCode >= 500;

  static NoteMeta _meta(NoteRow row) => NoteMeta(
    path: row.path,
    tags: row.tags,
    links: row.links,
    assignmentId: row.assignmentId,
    lessonId: row.lessonId,
  );
}
