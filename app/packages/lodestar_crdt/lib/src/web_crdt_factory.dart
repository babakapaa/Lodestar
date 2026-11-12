import 'dart:js_interop';

import 'crdt_text.dart';
import 'web_crdt_module.dart';
import 'web_crdt_text.dart';

const _modulePath = './crdt/lodestar_crdt.js';

WebCrdtModule? _module;

CrdtText createCrdtText(int clientId) {
  final module = _module;
  if (module == null) {
    throw StateError('CrdtText.ensureReady() was not awaited');
  }
  return WebCrdtText(module.create(clientId));
}

Future<void> prepareCrdt() async {
  if (_module != null) {
    return;
  }
  final module = await importModule(_modulePath.toJS).toDart as WebCrdtModule;
  await module.init().toDart;
  _module = module;
}
