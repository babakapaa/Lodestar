import '../../core/api/api_client.dart';
import '../../core/api/api_exception.dart';
import 'call.dart';
import 'call_ticket.dart';

class CallsApi {
  const CallsApi(this._api);

  static const notFound = 'call_not_found';

  final ApiClient _api;

  Future<Call?> active(String chatId) async {
    try {
      return Call.fromJson(
        (await _api.get(_call(chatId)))! as Map<String, Object?>,
      );
    } on ApiException catch (error) {
      if (error.code == notFound) {
        return null;
      }
      rethrow;
    }
  }

  Future<CallTicket> join(String chatId) async => CallTicket.fromJson(
    (await _api.post(_call(chatId)))! as Map<String, Object?>,
  );

  Future<void> leave(String callId) =>
      _api.post('/v1/chat/calls/${Uri.encodeComponent(callId)}/leave');

  static String _call(String chatId) =>
      '/v1/chat/chats/${Uri.encodeComponent(chatId)}/call';
}
