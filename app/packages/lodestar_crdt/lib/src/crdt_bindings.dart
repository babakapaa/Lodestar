import 'dart:ffi';

final class CrdtHandle extends Opaque {}

@Native<Pointer<Uint8> Function(Uint32)>(symbol: 'lc_alloc')
external Pointer<Uint8> crdtAlloc(int len);

@Native<Void Function(Pointer<Uint8>, Uint32)>(symbol: 'lc_dealloc')
external void crdtDealloc(Pointer<Uint8> ptr, int len);

@Native<Pointer<CrdtHandle> Function(Uint32)>(symbol: 'lc_doc_new')
external Pointer<CrdtHandle> crdtDocNew(int clientId);

@Native<Void Function(Pointer<CrdtHandle>)>(symbol: 'lc_doc_free')
external void crdtDocFree(Pointer<CrdtHandle> handle);

@Native<Pointer<Uint8> Function(Pointer<CrdtHandle>)>(symbol: 'lc_out_ptr')
external Pointer<Uint8> crdtOutPtr(Pointer<CrdtHandle> handle);

@Native<Int32 Function(Pointer<CrdtHandle>, Pointer<Uint8>, Uint32)>(
  symbol: 'lc_apply',
)
external int crdtApply(Pointer<CrdtHandle> handle, Pointer<Uint8> ptr, int len);

@Native<Uint32 Function(Pointer<CrdtHandle>)>(symbol: 'lc_state_vector')
external int crdtStateVector(Pointer<CrdtHandle> handle);

@Native<Int32 Function(Pointer<CrdtHandle>, Pointer<Uint8>, Uint32)>(
  symbol: 'lc_diff',
)
external int crdtDiff(Pointer<CrdtHandle> handle, Pointer<Uint8> ptr, int len);

@Native<Uint32 Function(Pointer<CrdtHandle>)>(symbol: 'lc_text')
external int crdtText(Pointer<CrdtHandle> handle);

@Native<Uint32 Function(Pointer<CrdtHandle>)>(symbol: 'lc_length')
external int crdtLength(Pointer<CrdtHandle> handle);

@Native<
  Int32 Function(Pointer<CrdtHandle>, Uint32, Uint32, Pointer<Uint8>, Uint32)
>(symbol: 'lc_edit')
external int crdtEdit(
  Pointer<CrdtHandle> handle,
  int index,
  int remove,
  Pointer<Uint8> ptr,
  int len,
);
