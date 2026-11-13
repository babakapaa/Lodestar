class DialogScope {
  const DialogScope({required this.kind, required this.id});

  final String kind;
  final String id;

  Map<String, Object?> toJson() => {'scopeKind': kind, 'scopeId': id};
}
