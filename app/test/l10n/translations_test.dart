import 'dart:convert';
import 'dart:io';

import 'package:flutter_test/flutter_test.dart';

Set<String> keysOf(String locale) {
  final file = File('lib/l10n/app_$locale.arb');
  final json = jsonDecode(file.readAsStringSync()) as Map<String, Object?>;
  return json.keys.where((key) => !key.startsWith('@')).toSet();
}

void main() {
  test('every language has every string', () {
    final reference = keysOf('ru');
    for (final locale in ['kk', 'en']) {
      expect(keysOf(locale), equals(reference), reason: 'app_$locale.arb');
    }
  });
}
