import 'dart:typed_data';

import 'crdt_factory.dart' if (dart.library.js_interop) 'web_crdt_factory.dart';

abstract interface class CrdtText {
  factory CrdtText(int clientId) => createCrdtText(clientId);

  static Future<void> ensureReady() => prepareCrdt();

  String get text;

  int get length;

  Uint8List get stateVector;

  void apply(Uint8List update);

  Uint8List diff(Uint8List remoteStateVector);

  void edit(int index, int remove, String insert);

  void dispose();
}
