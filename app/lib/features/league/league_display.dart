enum LeagueDisplay {
  name,
  hidden,
  nickname;

  static LeagueDisplay parse(Object? value) => LeagueDisplay.values.firstWhere(
    (d) => d.name == value,
    orElse: () => LeagueDisplay.name,
  );
}
