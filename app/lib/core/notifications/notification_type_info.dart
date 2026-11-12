import 'notification_text.dart';

class NotificationTypeInfo {
  const NotificationTypeInfo({
    required this.key,
    required this.required,
    required this.names,
    required this.texts,
  });

  factory NotificationTypeInfo.fromJson(Map<String, Object?> json) {
    return NotificationTypeInfo(
      key: json['key']! as String,
      required: json['required']! as bool,
      names: (json['names']! as Map<String, Object?>).map(
        (language, name) => MapEntry(language, name! as String),
      ),
      texts: (json['texts']! as Map<String, Object?>).map(
        (language, text) => MapEntry(
          language,
          NotificationText.fromJson(text! as Map<String, Object?>),
        ),
      ),
    );
  }

  static const _fallbackLanguage = 'en';

  final String key;
  final bool required;
  final Map<String, String> names;
  final Map<String, NotificationText> texts;

  String name(String language) => names[_pick(names.keys, language)] ?? key;

  NotificationText? render(Map<String, Object?> params, String language) {
    final text = texts[_pick(texts.keys, language)];
    if (text == null) {
      return null;
    }
    var title = text.title;
    var body = text.body;
    for (final MapEntry(:key, :value) in params.entries) {
      final placeholder = '{$key}';
      title = title.replaceAll(placeholder, '$value');
      body = body.replaceAll(placeholder, '$value');
    }
    return NotificationText(title: title, body: body);
  }

  static String? _pick(Iterable<String> languages, String preferred) {
    final available = languages.toSet();
    if (available.contains(preferred)) {
      return preferred;
    }
    if (available.contains(_fallbackLanguage)) {
      return _fallbackLanguage;
    }
    if (available.isEmpty) {
      return null;
    }
    return (available.toList()..sort()).first;
  }
}
