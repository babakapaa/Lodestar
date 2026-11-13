enum TableFormat {
  xlsx,
  csv;

  String get contentType => switch (this) {
    xlsx => 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    csv => 'text/csv',
  };

  static TableFormat ofName(String name) =>
      name.toLowerCase().endsWith('.xlsx') ? xlsx : csv;
}
