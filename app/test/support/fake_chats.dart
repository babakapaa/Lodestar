import 'dart:convert';

import 'package:http/http.dart' as http;

class FakeChats {
  FakeChats(this.files);

  static const me = 'u1';

  final Map<String, Map<String, Object?>> files;
  final Set<String> blocked = {};
  final List<Map<String, Object?>> reports = [];
  final Map<String, Map<String, String>> roles = {};
  final Map<String, Map<String, Object?>> calls = {};
  String callError = '';
  bool screenShare = true;
  bool record = false;

  final Map<String, Map<String, Object?>> chats = {};
  final Map<String, List<Map<String, Object?>>> messages = {};
  final Map<String, Map<String, String>> members = {};
  final List<Map<String, Object?>> sent = [];
  int _next = 1;

  Map<String, Object?> addChat({
    required String id,
    required String kind,
    required String title,
    String role = 'member',
    bool request = false,
    Map<String, String> people = const {},
    Map<String, Object?>? origin,
  }) {
    final chat = <String, Object?>{
      'id': id,
      'kind': kind,
      'title': title,
      'about': '',
      'origin': origin,
      'peerId': kind == 'direct' ? people.keys.first : null,
      'role': role,
      'muted': false,
      'request': request,
      'seq': 0,
      'rev': 0,
      'readSeq': 0,
      'unread': 0,
      'memberCount': people.length + 1,
      'last': null,
      'archivedAt': null,
      'createdAt': '2026-11-13T09:00:00Z',
      'changedAt': '2026-11-13T09:00:00Z',
    };
    chats[id] = chat;
    messages[id] = [];
    members[id] = {me: 'Айгерим Сапарова', ...people};
    roles[id] = {me: role};
    return chat;
  }

  Map<String, Object?> incoming(
    String chatId,
    String authorId,
    String text, {
    List<String> fileIds = const [],
  }) => _append(chatId, authorId, null, text, null, fileIds: fileIds);

  http.Response handle(http.Request request) {
    final segments = request.url.pathSegments.skip(2).toList();
    final query = request.url.queryParameters;
    return switch ((request.method, segments)) {
      ('GET', ['chats']) => _json(200, chats.values.toList()),
      ('POST', ['chats']) => _create(_body(request)),
      ('POST', ['direct']) => _direct(_body(request)),
      ('GET', ['chats', final id]) => _details(id),
      ('PUT', ['chats', final id, 'mute']) => _set(
        id,
        'muted',
        _body(request)['muted'],
      ),
      ('POST', ['chats', final id, 'accept']) => _set(id, 'request', false),
      ('POST', ['chats', final id, 'read']) => _read(id, _body(request)),
      ('GET', ['chats', final id, 'messages']) => _page(id, query),
      ('POST', ['chats', final id, 'messages']) => _send(id, _body(request)),
      ('PATCH', ['messages', final id]) => _change(
        id,
        _body(request)['text'] as String?,
      ),
      ('DELETE', ['messages', final id]) => _change(id, null),
      ('PATCH', ['chats', final id]) => _update(id, _body(request)),
      ('POST', ['chats', final id, 'members']) => _add(id, _body(request)),
      ('DELETE', ['chats', final id, 'members', final userId]) => _remove(
        id,
        userId,
      ),
      ('GET', ['chats', final id, 'call']) =>
        calls[id] == null
            ? _json(404, {'code': 'call_not_found'})
            : _json(200, calls[id]),
      ('POST', ['chats', final id, 'call']) => _joinCall(id),
      ('POST', ['calls', final id, 'leave']) => _leaveCall(id),
      ('POST', ['calls', final id, 'recording']) => _recording(id, true),
      ('DELETE', ['calls', final id, 'recording']) => _recording(id, false),
      ('GET', ['blocks']) => _json(200, [
        for (final id in blocked) {'userId': id, 'displayName': id},
      ]),
      ('POST', ['blocks']) => _block(_body(request)['userId']! as String),
      ('DELETE', ['blocks', final userId]) => _unblock(userId),
      ('POST', ['reports']) => _report(_body(request)),
      ('GET', ['reports']) => _json(200, [
        for (final r in reports)
          if (r['resolvedAt'] == null) r,
      ]),
      ('POST', ['reports', final id, 'resolve']) => _resolve(id),
      _ => _json(404, {'code': 'not_found'}),
    };
  }

  Map<String, Object?> startCall(String chatId, List<String> people) {
    final call = <String, Object?>{
      'id': 'call-$chatId',
      'chatId': chatId,
      'startedBy': people.isEmpty ? null : people.first,
      'startedAt': '2026-11-13T10:00:00Z',
      'people': people,
      'recording': false,
    };
    calls[chatId] = call;
    return call;
  }

  http.Response _joinCall(String chatId) {
    if (callError.isNotEmpty) {
      return _json(403, {'code': callError});
    }
    final call = calls[chatId] ?? startCall(chatId, const []);
    call['people'] = {...call['people']! as List<String>, me}.toList();
    return _json(200, {
      'call': call,
      'url': 'wss://calls.example.edu',
      'room': 'kaznu_${call['id']}',
      'token': 'room-token',
      'screenShare': screenShare,
      'record': record,
      'muted':
          chats[chatId]?['kind'] == 'channel' &&
          chats[chatId]?['role'] == 'member',
    });
  }

  http.Response _recording(String callId, bool on) {
    final call = calls.values.firstWhere((call) => call['id'] == callId);
    if (!record) {
      return _json(403, {'code': 'recording_off'});
    }
    call['recording'] = on;
    return _json(200, call);
  }

  http.Response _leaveCall(String callId) {
    calls.removeWhere((_, call) => call['id'] == callId);
    return http.Response('', 204);
  }

  http.Response _create(Map<String, Object?> body) {
    final id = 'c${_next++}';
    addChat(
      id: id,
      kind: body['kind']! as String,
      title: (body['title'] as String?) ?? '',
      role: 'owner',
      origin: body['origin'] as Map<String, Object?>?,
    );
    return _json(201, _full(id));
  }

  http.Response _direct(Map<String, Object?> body) {
    final userId = body['userId']! as String;
    for (final chat in chats.values) {
      if (chat['kind'] == 'direct' && chat['peerId'] == userId) {
        return _json(200, _full(chat['id']! as String));
      }
    }
    final id = 'd${_next++}';
    addChat(
      id: id,
      kind: 'direct',
      title: 'Ерлан Ким',
      people: {userId: 'Ерлан Ким'},
    );
    return _json(201, _full(id));
  }

  http.Response _details(String id) => chats.containsKey(id)
      ? _json(200, _full(id))
      : _json(404, {'code': 'chat_not_found'});

  http.Response _set(String id, String key, Object? value) {
    chats[id]![key] = value;
    return _json(200, chats[id]);
  }

  http.Response _read(String id, Map<String, Object?> body) {
    final chat = chats[id]!;
    final seq = body['seq']! as int;
    chat['readSeq'] = seq;
    chat['unread'] = (chat['seq']! as int) - seq;
    return _json(200, chat);
  }

  http.Response _page(String id, Map<String, String> query) {
    final all = messages[id]!;
    final since = int.tryParse(query['since'] ?? '');
    final before = int.tryParse(query['before'] ?? '');
    final list = since != null
        ? ([
            for (final m in all)
              if ((m['rev']! as int) > since) m,
          ]..sort((a, b) => (a['rev']! as int).compareTo(b['rev']! as int)))
        : [
            for (final m in all)
              if (before == null || (m['seq']! as int) < before) m,
          ];
    return _json(200, {'messages': list, 'more': false});
  }

  http.Response _send(String id, Map<String, Object?> body) {
    final clientId = body['clientId']! as String;
    for (final m in messages[id]!) {
      if (m['clientId'] == clientId) {
        return _json(200, m);
      }
    }
    sent.add(body);
    final m = _append(
      id,
      me,
      clientId,
      body['text'] as String? ?? '',
      body['replyTo'] as int?,
      fileIds: [
        for (final id in body['fileIds'] as List<Object?>? ?? const [])
          id! as String,
      ],
      mentions: [
        for (final id in body['mentions'] as List<Object?>? ?? const [])
          id! as String,
      ],
    );
    chats[id]!['readSeq'] = m['seq'];
    chats[id]!['unread'] = 0;
    return _json(201, m);
  }

  http.Response _change(String id, String? text) {
    for (final entry in messages.entries) {
      for (final m in entry.value) {
        if (m['id'] == id) {
          final chat = chats[entry.key]!;
          chat['rev'] = (chat['rev']! as int) + 1;
          m['rev'] = chat['rev'];
          if (text == null) {
            m['deleted'] = true;
            m['text'] = '';
          } else {
            m['text'] = text;
            m['editedAt'] = '2026-11-13T10:30:00Z';
          }
          return _json(200, m);
        }
      }
    }
    return _json(404, {'code': 'message_not_found'});
  }

  Map<String, Object?> _append(
    String chatId,
    String authorId,
    String? clientId,
    String text,
    int? replyTo, {
    List<String> fileIds = const [],
    List<String> mentions = const [],
  }) {
    final chat = chats[chatId]!;
    final seq = (chat['seq']! as int) + 1;
    final rev = (chat['rev']! as int) + 1;
    final m = <String, Object?>{
      'id': 'm${_next++}',
      'chatId': chatId,
      'seq': seq,
      'rev': rev,
      'authorId': authorId,
      'clientId': clientId,
      'text': text,
      'replyTo': replyTo,
      'mentions': mentions,
      'files': [for (final id in fileIds) files[id]],
      'views': chat['kind'] == 'channel' ? 0 : null,
      'deleted': false,
      'editedAt': null,
      'createdAt': '2026-11-13T10:0${seq % 10}:00Z',
    };
    messages[chatId]!.add(m);
    chat['seq'] = seq;
    chat['rev'] = rev;
    chat['last'] = m;
    chat['changedAt'] = m['createdAt'];
    chat['unread'] = seq - (chat['readSeq']! as int);
    return m;
  }

  Map<String, Object?> _full(String id) => {
    ...chats[id]!,
    'members': [
      for (final MapEntry(:key, :value) in members[id]!.entries)
        {
          'userId': key,
          'displayName': value,
          'role': roles[id]![key] ?? 'member',
          'synced': false,
          'joinedAt': '2026-11-13T09:00:00Z',
        },
    ],
    'peerReadSeq': null,
  };

  http.Response _update(String id, Map<String, Object?> body) {
    chats[id]!['title'] = body['title'];
    chats[id]!['about'] = body['about'];
    return _json(200, _full(id));
  }

  http.Response _add(String id, Map<String, Object?> body) {
    for (final userId in body['userIds']! as List<Object?>) {
      members[id]![userId! as String] = 'Ерлан Ким';
    }
    chats[id]!['memberCount'] = members[id]!.length;
    return _json(200, _full(id));
  }

  http.Response _remove(String id, String userId) {
    members[id]!.remove(userId);
    if (userId == me) {
      chats.remove(id);
    } else {
      chats[id]!['memberCount'] = members[id]!.length;
    }
    return http.Response('', 204);
  }

  http.Response _block(String userId) {
    blocked.add(userId);
    return http.Response('', 204);
  }

  http.Response _unblock(String userId) {
    blocked.remove(userId);
    return http.Response('', 204);
  }

  http.Response _report(Map<String, Object?> body) {
    final chatId = body['chatId']! as String;
    final ids = body['messageIds']! as List<Object?>;
    final report = <String, Object?>{
      'id': 'r${_next++}',
      'chatId': chatId,
      'chatTitle': chats[chatId]!['title'],
      'reporterId': me,
      'reason': body['reason'],
      'messages': [
        for (final m in messages[chatId]!)
          if (ids.contains(m['id']))
            {
              'seq': m['seq'],
              'authorId': m['authorId'],
              'authorName': members[chatId]![m['authorId']] ?? '',
              'text': m['text'],
              'createdAt': m['createdAt'],
            },
      ],
      'createdAt': '2026-11-13T11:00:00Z',
      'resolvedAt': null,
    };
    reports.add(report);
    return _json(201, report);
  }

  http.Response _resolve(String id) {
    final report = reports.firstWhere((r) => r['id'] == id);
    report['resolvedAt'] = '2026-11-13T12:00:00Z';
    return _json(200, report);
  }

  static Map<String, Object?> _body(http.Request request) =>
      jsonDecode(request.body) as Map<String, Object?>;

  static http.Response _json(int status, Object? body) => http.Response(
    jsonEncode(body),
    status,
    headers: {'content-type': 'application/json; charset=utf-8'},
  );
}
