class Contact {
  const Contact({
    required this.kind,
    required this.value,
    required this.public,
  });

  factory Contact.fromJson(Map<String, Object?> json) {
    return Contact(
      kind: json['kind']! as String,
      value: json['value']! as String,
      public: json['public']! as bool,
    );
  }

  final String kind;
  final String value;
  final bool public;

  Map<String, Object?> toJson() => {
    'kind': kind,
    'value': value,
    'public': public,
  };
}
