import 'dart:async';
import 'dart:math';
import 'dart:typed_data';

import 'package:lodestar_crdt/lodestar_crdt.dart';

import 'text_edit.dart';

class NoteDoc {
  NoteDoc(this.id, Uint8List state) : _text = CrdtText(_clientId()) {
    if (state.isNotEmpty) {
      _text.apply(state);
    }
    _flushed = _text.stateVector;
  }

  static final _random = Random.secure();
  static final _empty = Uint8List.fromList([0]);

  final String id;
  final CrdtText _text;
  final _remote = StreamController<TextEdit>.broadcast();
  late Uint8List _flushed;
  bool _dirty = false;

  String get text => _text.text;

  bool get dirty => _dirty;

  Uint8List get state => _text.diff(_empty);

  Stream<TextEdit> get remoteEdits => _remote.stream;

  bool edit(String next) {
    final change = TextEdit.between(_text.text, next);
    if (change == null) {
      return false;
    }
    _text.edit(change.index, change.remove, change.insert);
    _dirty = true;
    return true;
  }

  Uint8List? takeLocal() {
    if (!_dirty) {
      return null;
    }
    final update = _text.diff(_flushed);
    _flushed = _text.stateVector;
    _dirty = false;
    return update;
  }

  void applyRemote(List<Uint8List> updates) {
    final before = _text.text;
    for (final update in updates) {
      _text.apply(update);
    }
    if (!_dirty) {
      _flushed = _text.stateVector;
    }
    final change = TextEdit.between(before, _text.text);
    if (change != null) {
      _remote.add(change);
    }
  }

  void dispose() {
    _remote.close();
    _text.dispose();
  }

  static int _clientId() => _random.nextInt(0xFFFFFFFF) + 1;
}
