class SetupSection {
  const SetupSection({
    required this.key,
    required this.required,
    required this.status,
    required this.assigneeId,
    required this.dueOn,
    required this.note,
    required this.settings,
  });

  factory SetupSection.fromJson(Map<String, Object?> json) => SetupSection(
    key: json['key']! as String,
    required: json['required'] as bool? ?? false,
    status: json['status'] as String? ?? notStarted,
    assigneeId: json['assigneeId'] as String?,
    dueOn: json['dueOn'] as String?,
    note: json['note'] as String? ?? '',
    settings: [
      for (final item in json['settings'] as List<Object?>? ?? const [])
        item! as String,
    ],
  );

  static const notStarted = 'not_started';
  static const inProgress = 'in_progress';
  static const review = 'review';
  static const approved = 'approved';
  static const legal = 'legal';
  static const deploy = 'deploy';

  final String key;
  final bool required;
  final String status;
  final String? assigneeId;
  final String? dueOn;
  final String note;
  final List<String> settings;

  bool get isApproved => status == approved;
}
