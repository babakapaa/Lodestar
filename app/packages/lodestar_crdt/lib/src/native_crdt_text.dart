import 'dart:convert';
import 'dart:ffi';
import 'dart:typed_data';

import 'crdt_bindings.dart';
import 'crdt_exception.dart';
import 'crdt_text.dart';

final class NativeCrdtText implements CrdtText, Finalizable {
  NativeCrdtText(int clientId) : _handle = crdtDocNew(clientId) {
    _finalizer.attach(this, _handle.cast(), detach: this);
  }

  static final _finalizer = NativeFinalizer(
    Native.addressOf<NativeFunction<Void Function(Pointer<CrdtHandle>)>>(
      crdtDocFree,
    ).cast(),
  );

  Pointer<CrdtHandle> _handle;

  @override
  String get text => utf8.decode(_out(crdtText(_live)));

  @override
  int get length => crdtLength(_live);

  @override
  Uint8List get stateVector => _out(crdtStateVector(_live));

  @override
  void apply(Uint8List update) {
    _check(_withBytes(update, (ptr, len) => crdtApply(_live, ptr, len)));
  }

  @override
  Uint8List diff(Uint8List remoteStateVector) {
    final len = _withBytes(
      remoteStateVector,
      (ptr, len) => crdtDiff(_live, ptr, len),
    );
    if (len < 0) {
      throw CrdtException(-len);
    }
    return _out(len);
  }

  @override
  void edit(int index, int remove, String insert) {
    final bytes = utf8.encode(insert);
    _check(
      _withBytes(bytes, (ptr, len) => crdtEdit(_live, index, remove, ptr, len)),
    );
  }

  @override
  void dispose() {
    if (_handle == nullptr) {
      return;
    }
    _finalizer.detach(this);
    crdtDocFree(_handle);
    _handle = nullptr;
  }

  Pointer<CrdtHandle> get _live {
    if (_handle == nullptr) {
      throw StateError('CrdtText is disposed');
    }
    return _handle;
  }

  Uint8List _out(int len) =>
      Uint8List.fromList(crdtOutPtr(_live).asTypedList(len));

  static int _withBytes(
    Uint8List bytes,
    int Function(Pointer<Uint8> ptr, int len) call,
  ) {
    final ptr = crdtAlloc(bytes.length);
    try {
      ptr.asTypedList(bytes.length).setAll(0, bytes);
      return call(ptr, bytes.length);
    } finally {
      crdtDealloc(ptr, bytes.length);
    }
  }

  static void _check(int code) {
    if (code != 0) {
      throw CrdtException(code);
    }
  }
}
