enum VoteDecision {
  waiting,
  agreed,
  objected,
  silent;

  static VoteDecision parse(Object? value) => VoteDecision.values.firstWhere(
    (d) => d.name == value,
    orElse: () => VoteDecision.waiting,
  );
}
