import 'package:drift/drift.dart';

import 'note_rows.dart';
import 'pending_updates.dart';
import 'string_list_converter.dart';

part 'notes_database.g.dart';

@DriftDatabase(tables: [NoteRows, PendingUpdates])
class NotesDatabase extends _$NotesDatabase {
  NotesDatabase(super.executor);

  @override
  int get schemaVersion => 1;

  Future<void> clear() => transaction(() async {
    await delete(pendingUpdates).go();
    await delete(noteRows).go();
  });
}
