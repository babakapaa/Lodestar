import 'dart:js_interop';

extension type WebCrdtDoc._(JSObject _) implements JSObject {
  external int apply(JSUint8Array update);

  external JSUint8Array stateVector();

  external JSUint8Array? diff(JSUint8Array remote);

  external String text();

  external int length();

  external int edit(int index, int remove, String insert);

  external void free();
}
