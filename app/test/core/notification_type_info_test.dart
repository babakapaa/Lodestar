import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/notifications/notification_text.dart';
import 'package:lodestar/core/notifications/notification_type_info.dart';

void main() {
  const type = NotificationTypeInfo(
    key: 'study.deadline',
    required: false,
    names: {'ru': 'Сроки', 'en': 'Deadlines'},
    texts: {
      'ru': NotificationText(title: 'Срок: {course}', body: '{days} дн.'),
      'en': NotificationText(title: 'Due: {course}', body: '{days} days'),
    },
  );

  test('uses the app language and fills params', () {
    final text = type.render({'course': 'Химия', 'days': 2}, 'ru')!;

    expect(text.title, 'Срок: Химия');
    expect(text.body, '2 дн.');
  });

  test('falls back to English, then to any language', () {
    expect(type.name('kk'), 'Deadlines');
    const only = NotificationTypeInfo(
      key: 'a.b',
      required: false,
      names: {'ru': 'Только', 'kk': 'Тек'},
      texts: {},
    );
    expect(only.name('en'), 'Тек');
    expect(only.render(const {}, 'en'), isNull);
  });
}
