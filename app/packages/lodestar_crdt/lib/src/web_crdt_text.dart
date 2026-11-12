import 'dart:js_interop';
import 'dart:typed_data';

import 'crdt_exception.dart';
import 'crdt_text.dart';
import 'web_crdt_doc.dart';

final class WebCrdtText implements CrdtText {
  WebCrdtText(this._doc);

  WebCrdtDoc? _doc;

  @override
  String get text => _live.text();

  @override
  int get length => _live.length();

  @override
  Uint8List get stateVector => _live.stateVector().toDart;

  @override
  void apply(Uint8List update) => _check(_live.apply(update.toJS));

  @override
  Uint8List diff(Uint8List remoteStateVector) {
    final out = _live.diff(remoteStateVector.toJS);
    if (out == null) {
      throw const CrdtException(CrdtException.badInput);
    }
    return out.toDart;
  }

  @override
  void edit(int index, int remove, String insert) =>
      _check(_live.edit(index, remove, insert));

  @override
  void dispose() {
    _doc?.free();
    _doc = null;
  }

  WebCrdtDoc get _live {
    final doc = _doc;
    if (doc == null) {
      throw StateError('CrdtText is disposed');
    }
    return doc;
  }

  static void _check(int code) {
    if (code != 0) {
      throw CrdtException(code);
    }
  }
}
