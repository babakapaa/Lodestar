import 'package:drift/drift.dart';

class PendingUpdates extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get noteId => text()();
  TextColumn get key => text()();
  BlobColumn get data => blob()();
}
