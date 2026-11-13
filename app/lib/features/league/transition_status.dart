enum TransitionStatus {
  pending,
  applied,
  rejected,
  reverted;

  static TransitionStatus parse(Object? value) =>
      TransitionStatus.values.firstWhere(
        (s) => s.name == value,
        orElse: () => TransitionStatus.pending,
      );
}
