class Deletion {
  const Deletion({required this.dueAt});

  factory Deletion.fromJson(Map<String, Object?> json) {
    return Deletion(dueAt: DateTime.parse(json['dueAt']! as String));
  }

  final DateTime dueAt;
}
