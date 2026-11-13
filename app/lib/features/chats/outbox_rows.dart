import 'package:drift/drift.dart';

class OutboxRows extends Table {
  TextColumn get clientId => text()();
  TextColumn get chatId => text()();
  TextColumn get body => text()();
  IntColumn get replyTo => integer().nullable()();
  DateTimeColumn get createdAt => dateTime()();
  TextColumn get error => text().nullable()();

  @override
  Set<Column<Object>> get primaryKey => {clientId};
}
