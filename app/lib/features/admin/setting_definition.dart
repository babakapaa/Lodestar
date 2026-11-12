class SettingDefinition {
  const SettingDefinition({
    required this.key,
    required this.type,
    required this.options,
    required this.min,
    required this.max,
    required this.maxLength,
    required this.defaultValue,
  });

  factory SettingDefinition.fromJson(Map<String, Object?> json) {
    final options = json['options'];
    return SettingDefinition(
      key: json['key']! as String,
      type: json['type']! as String,
      options: [
        if (options is List<Object?>)
          for (final option in options) option! as String,
      ],
      min: (json['min'] as num?)?.toInt(),
      max: (json['max'] as num?)?.toInt(),
      maxLength: (json['maxLength'] as num?)?.toInt(),
      defaultValue: json['default'],
    );
  }

  static const boolType = 'bool';
  static const integerType = 'integer';
  static const textType = 'text';
  static const choiceType = 'choice';
  static const timeZoneType = 'time_zone';

  final String key;
  final String type;
  final List<String> options;
  final int? min;
  final int? max;
  final int? maxLength;
  final Object? defaultValue;
}
