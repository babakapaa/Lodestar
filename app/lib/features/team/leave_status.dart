enum LeaveStatus {
  pending,
  approved,
  rejected;

  static LeaveStatus parse(Object? value) => LeaveStatus.values.firstWhere(
    (s) => s.name == value,
    orElse: () => LeaveStatus.pending,
  );
}
