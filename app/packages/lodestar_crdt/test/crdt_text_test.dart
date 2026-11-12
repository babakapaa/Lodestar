import 'dart:typed_data';

import 'package:lodestar_crdt/lodestar_crdt.dart';
import 'package:test/test.dart';

void sync(CrdtText from, CrdtText to) => to.apply(from.diff(to.stateVector));

void main() {
  setUpAll(CrdtText.ensureReady);

  test('offline edits from two people merge without loss', () {
    final a = CrdtText(1);
    final b = CrdtText(2);
    a.edit(0, 0, 'Клетка состоит из ядра.');
    sync(a, b);
    a.edit(0, 0, 'Айгерим: ');
    b.edit(b.length, 0, ' И мембраны 🧬.');
    sync(b, a);
    sync(a, b);
    expect(a.text, 'Айгерим: Клетка состоит из ядра. И мембраны 🧬.');
    expect(b.text, a.text);
    a.dispose();
    b.dispose();
  });

  test('offsets are UTF-16 like Flutter text fields', () {
    final doc = CrdtText(1)..edit(0, 0, 'a🧬b');
    expect(doc.length, 'a🧬b'.length);
    doc.edit(1, 2, '');
    expect(doc.text, 'ab');
  });

  test('bad input throws instead of crashing', () {
    final doc = CrdtText(1);
    expect(() => doc.edit(5, 0, 'x'), throwsA(isA<CrdtException>()));
    expect(
      () => doc.apply(Uint8List.fromList([255, 255, 255])),
      throwsA(isA<CrdtException>()),
    );
    doc.dispose();
    expect(() => doc.text, throwsStateError);
  });
}
