import '../../core/api/api_client.dart';
import '../../core/offline/cached_reader.dart';
import 'blocked_user.dart';
import 'chat.dart';
import 'chat_details.dart';
import 'chat_message.dart';
import 'chat_report.dart';
import 'message_page.dart';
import 'new_chat.dart';
import 'outgoing.dart';

class ChatsApi {
  const ChatsApi(this._api, this._reader);

  final ApiClient _api;
  final CachedReader _reader;

  Future<String> me() async =>
      ((await _api.get('/v1/me'))! as Map<String, Object?>)['id']! as String;

  Future<List<Chat>> list() async => [
    for (final item in (await _api.get('/v1/chat/chats'))! as List<Object?>)
      Chat.fromJson(item! as Map<String, Object?>),
  ];

  Future<ChatDetails> get(String id) async => ChatDetails.fromJson(
    (await _reader.get(_chat(id))).body! as Map<String, Object?>,
  );

  Future<ChatDetails> create(NewChat chat) async => ChatDetails.fromJson(
    (await _api.post('/v1/chat/chats', chat.toJson()))! as Map<String, Object?>,
  );

  Future<ChatDetails> direct(String userId) async => ChatDetails.fromJson(
    (await _api.post('/v1/chat/direct', {'userId': userId}))!
        as Map<String, Object?>,
  );

  Future<ChatDetails> update(String id, String title, String about) async =>
      ChatDetails.fromJson(
        (await _api.patch(_chat(id), {'title': title, 'about': about}))!
            as Map<String, Object?>,
      );

  Future<ChatDetails> addMembers(String id, List<String> userIds) async =>
      ChatDetails.fromJson(
        (await _api.post('${_chat(id)}/members', {'userIds': userIds}))!
            as Map<String, Object?>,
      );

  Future<void> removeMember(String id, String userId) =>
      _api.delete('${_chat(id)}/members/${Uri.encodeComponent(userId)}');

  Future<Chat> mute(String id, {required bool muted}) async => Chat.fromJson(
    (await _api.put('${_chat(id)}/mute', {'muted': muted}))!
        as Map<String, Object?>,
  );

  Future<Chat> accept(String id) async => Chat.fromJson(
    (await _api.post('${_chat(id)}/accept'))! as Map<String, Object?>,
  );

  Future<Chat> read(String id, int seq) async => Chat.fromJson(
    (await _api.post('${_chat(id)}/read', {'seq': seq}))!
        as Map<String, Object?>,
  );

  Future<MessagePage> latest(String id) => _page(id, const {});

  Future<MessagePage> before(String id, int seq) =>
      _page(id, {'before': '$seq'});

  Future<MessagePage> since(String id, int rev) =>
      _page(id, {'since': '$rev', 'limit': '200'});

  Future<ChatMessage> send(Outgoing message) async => ChatMessage.fromJson(
    (await _api.post('${_chat(message.chatId)}/messages', message.toJson()))!
        as Map<String, Object?>,
  );

  Future<ChatMessage> edit(String id, String text) async =>
      ChatMessage.fromJson(
        (await _api.patch(_message(id), {'text': text}))!
            as Map<String, Object?>,
      );

  Future<ChatMessage> delete(String id) async => ChatMessage.fromJson(
    (await _api.delete(_message(id)))! as Map<String, Object?>,
  );

  Future<List<BlockedUser>> blocks() async => [
    for (final item in (await _api.get('/v1/chat/blocks'))! as List<Object?>)
      BlockedUser.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> block(String userId) =>
      _api.post('/v1/chat/blocks', {'userId': userId});

  Future<void> unblock(String userId) =>
      _api.delete('/v1/chat/blocks/${Uri.encodeComponent(userId)}');

  Future<void> report(String chatId, List<String> messageIds, String reason) =>
      _api.post('/v1/chat/reports', {
        'chatId': chatId,
        'messageIds': messageIds,
        'reason': reason,
      });

  Future<List<ChatReport>> reports() async => [
    for (final item in (await _api.get('/v1/chat/reports'))! as List<Object?>)
      ChatReport.fromJson(item! as Map<String, Object?>),
  ];

  Future<ChatReport> resolve(String id) async => ChatReport.fromJson(
    (await _api.post('/v1/chat/reports/${Uri.encodeComponent(id)}/resolve'))!
        as Map<String, Object?>,
  );

  Future<MessagePage> _page(String id, Map<String, String> query) async =>
      MessagePage.fromJson(
        (await _api.get('${_chat(id)}/messages', query))!
            as Map<String, Object?>,
      );

  static String _chat(String id) => '/v1/chat/chats/${Uri.encodeComponent(id)}';

  static String _message(String id) =>
      '/v1/chat/messages/${Uri.encodeComponent(id)}';
}
