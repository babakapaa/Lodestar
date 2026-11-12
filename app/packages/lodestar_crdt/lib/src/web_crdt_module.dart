import 'dart:js_interop';
import 'dart:js_interop_unsafe';

import 'web_crdt_doc.dart';

extension type WebCrdtModule._(JSObject _) implements JSObject {
  @JS('default')
  external JSPromise<JSAny?> init();

  WebCrdtDoc create(int clientId) => (getProperty<JSFunction>(
    'CrdtDoc'.toJS,
  ).callAsConstructor<WebCrdtDoc>(clientId.toJS));
}
