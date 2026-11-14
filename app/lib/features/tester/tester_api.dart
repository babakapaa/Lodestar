import '../../core/api/api_exception.dart';
import '../../core/api/api_transport.dart';
import '../../core/api/server_address.dart';
import '../../core/api/token_pair.dart';
import 'tester_info.dart';
import 'tester_person.dart';

class TesterApi {
  const TesterApi(this._transport);

  final ApiTransport _transport;

  Future<TesterInfo?> info(ServerAddress address) async {
    try {
      final body = await _transport.send('GET', address.resolve('/v1/tester'));
      final info = TesterInfo.fromJson(body! as Map<String, Object?>);
      return info.tenant == address.tenant ? info : null;
    } on ApiException {
      return null;
    }
  }

  Future<TokenPair> signIn(
    ServerAddress address,
    TesterPerson person, {
    required String device,
  }) async {
    final body = await _transport.send(
      'POST',
      address.resolve('/v1/tester/sign-in'),
      body: {'person': person.key, 'device': device},
    );
    return TokenPair.fromJson(body! as Map<String, Object?>);
  }
}
