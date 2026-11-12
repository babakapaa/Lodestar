import 'package:drift/drift.dart';

import 'string_list_converter.dart';

class NoteRows extends Table {
  TextColumn get id => text()();
  TextColumn get ownerId => text().withDefault(const Constant(''))();
  TextColumn get role => text()();
  TextColumn get path => text()();
  TextColumn get tags => text().map(const StringListConverter())();
  TextColumn get links => text().map(const StringListConverter())();
  TextColumn get assignmentId => text().nullable()();
  TextColumn get lessonId => text().nullable()();
  IntColumn get seq => integer().withDefault(const Constant(0))();
  IntColumn get unsnapshotted => integer().withDefault(const Constant(0))();
  BlobColumn get state => blob()();
  TextColumn get body => text().withDefault(const Constant(''))();
  DateTimeColumn get changedAt => dateTime()();
  DateTimeColumn get deletedAt => dateTime().nullable()();
  BoolColumn get created => boolean().withDefault(const Constant(false))();
  BoolColumn get metaDirty => boolean().withDefault(const Constant(false))();
  TextColumn get error => text().nullable()();

  @override
  Set<Column<Object>> get primaryKey => {id};
}
