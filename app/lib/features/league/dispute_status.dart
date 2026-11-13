enum DisputeStatus {
  open,
  accepted,
  rejected;

  static DisputeStatus parse(Object? value) => DisputeStatus.values.firstWhere(
    (s) => s.name == value,
    orElse: () => DisputeStatus.open,
  );
}
