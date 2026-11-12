import 'dart:async';

import 'package:web_socket_channel/web_socket_channel.dart';

import 'live_socket.dart';

class WebLiveSocket implements LiveSocket {
  const WebLiveSocket();

  @override
  Stream<String> connect(Uri uri) {
    WebSocketChannel? channel;
    late final StreamController<String> controller;
    controller = StreamController<String>(
      onListen: () {
        final opened = WebSocketChannel.connect(uri);
        channel = opened;
        controller
            .addStream(opened.stream.where((m) => m is String).cast<String>())
            .whenComplete(controller.close);
      },
      onCancel: () => channel?.sink.close(),
    );
    return controller.stream;
  }
}
