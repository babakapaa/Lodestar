class NavigatorDialog {
  const NavigatorDialog({
    required this.id,
    required this.scopeKind,
    required this.title,
    required this.changedAt,
  });

  factory NavigatorDialog.fromJson(Map<String, Object?> json) =>
      NavigatorDialog(
        id: json['id']! as String,
        scopeKind: json['scopeKind'] as String? ?? '',
        title: json['title'] as String? ?? '',
        changedAt: DateTime.parse(json['changedAt']! as String),
      );

  final String id;
  final String scopeKind;
  final String title;
  final DateTime changedAt;
}
