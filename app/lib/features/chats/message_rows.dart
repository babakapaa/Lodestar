import 'package:drift/drift.dart';

class MessageRows extends Table {
  TextColumn get id => text()();
  TextColumn get chatId => text()();
  IntColumn get seq => integer()();
  TextColumn get data => text()();

  @override
  Set<Column<Object>> get primaryKey => {id};
}
