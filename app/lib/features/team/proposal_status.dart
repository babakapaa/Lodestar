enum ProposalStatus {
  open,
  sent,
  replaced;

  static ProposalStatus parse(Object? value) => ProposalStatus.values
      .firstWhere((s) => s.name == value, orElse: () => ProposalStatus.open);
}
