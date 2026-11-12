import 'dart:async';
import 'dart:convert';

import 'package:lodestar/core/live/live_socket.dart';

class FakeLiveSocket implements LiveSocket {
  final List<Uri> opened = [];
  StreamController<String>? _current;

  @override
  Stream<String> connect(Uri uri) {
    opened.add(uri);
    final controller = StreamController<String>();
    _current = controller;
    return controller.stream;
  }

  void send(int id, String type, Object? payload) {
    _current!.add(
      jsonEncode({
        'id': id,
        'type': type,
        'payload': payload,
        'at': '2026-11-12T10:00:00Z',
      }),
    );
  }

  Future<void> drop() async => _current!.close();
}
