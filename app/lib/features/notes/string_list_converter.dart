import 'dart:convert';

import 'package:drift/drift.dart';

class StringListConverter extends TypeConverter<List<String>, String> {
  const StringListConverter();

  @override
  List<String> fromSql(String fromDb) => [
    for (final item in jsonDecode(fromDb) as List<Object?>) item! as String,
  ];

  @override
  String toSql(List<String> value) => jsonEncode(value);
}
