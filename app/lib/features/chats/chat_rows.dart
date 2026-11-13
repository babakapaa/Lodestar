import 'package:drift/drift.dart';

class ChatRows extends Table {
  TextColumn get id => text()();
  TextColumn get data => text()();
  DateTimeColumn get changedAt => dateTime()();
  IntColumn get syncedRev => integer().withDefault(const Constant(-1))();
  BoolColumn get older => boolean().withDefault(const Constant(true))();

  @override
  Set<Column<Object>> get primaryKey => {id};
}
