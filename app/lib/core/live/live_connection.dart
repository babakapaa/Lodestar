import 'dart:async';
import 'dart:convert';
import 'dart:math';

import '../api/server_address.dart';
import '../api/token_store.dart';
import 'live_event.dart';
import 'live_socket.dart';

class LiveConnection {
  LiveConnection({
    required this._socket,
    required this._address,
    required this._tokens,
    required this._onFailure,
  });

  static const _maxDelaySeconds = 60;

  final LiveSocket _socket;
  final ServerAddress _address;
  final TokenStore _tokens;
  final Future<void> Function() _onFailure;
  final _events = StreamController<LiveEvent>.broadcast();
  StreamSubscription<String>? _subscription;
  Timer? _timer;
  int? _after;
  int _attempt = 0;
  bool _closed = false;

  Stream<LiveEvent> get events => _events.stream;

  void start() => _connect();

  Uri _uri(String accessToken) {
    final base = _address.resolve('/v1/live', {
      'access_token': accessToken,
      if (_after case final after?) 'after': '$after',
    });
    return base.replace(scheme: base.scheme == 'https' ? 'wss' : 'ws');
  }

  void _connect() {
    final tokens = _tokens.current;
    if (_closed || tokens == null) {
      return;
    }
    _subscription = _socket
        .connect(_uri(tokens.accessToken))
        .listen(
          _receive,
          onError: (Object _) => _retry(failed: true),
          onDone: _retry,
          cancelOnError: true,
        );
  }

  void _receive(String message) {
    _attempt = 0;
    final LiveEvent event;
    try {
      event = LiveEvent.fromJson(jsonDecode(message) as Map<String, Object?>);
    } on Object {
      return;
    }
    if (event.type == LiveEvent.resync) {
      final payload = event.payload;
      final latest = payload is Map<String, Object?> ? payload['latest'] : null;
      _after = latest is num ? latest.toInt() : null;
    } else {
      _after = event.id;
    }
    _events.add(event);
  }

  void _retry({bool failed = false}) {
    _subscription = null;
    if (_closed) {
      return;
    }
    if (failed) {
      unawaited(_onFailure());
    }
    final seconds = min(_maxDelaySeconds, pow(2, _attempt).toInt());
    _attempt++;
    _timer = Timer(Duration(seconds: seconds), _connect);
  }

  Future<void> close() async {
    _closed = true;
    _timer?.cancel();
    await _subscription?.cancel();
    await _events.close();
  }
}
