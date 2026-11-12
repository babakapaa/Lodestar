import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/api/server_address.dart';

void main() {
  test('adds https and trims the trailing slash', () {
    final address = ServerAddress.tryParse(
      ' campus.example.edu/lodestar/ ',
      'KazNU',
    )!;

    expect(address.server.toString(), 'https://campus.example.edu/lodestar');
    expect(address.tenant, 'kaznu');
    expect(
      address.resolve('/v1/me').toString(),
      'https://campus.example.edu/lodestar/v1/me',
    );
  });

  test('keeps an explicit scheme and port', () {
    final address = ServerAddress.tryParse('http://10.0.2.2:8080', 'demo')!;

    expect(address.server.toString(), 'http://10.0.2.2:8080');
  });

  test('rejects bad input', () {
    expect(ServerAddress.tryParse('', 'kaznu'), isNull);
    expect(ServerAddress.tryParse('ftp://host', 'kaznu'), isNull);
    expect(ServerAddress.tryParse('host', 'Bad code!'), isNull);
    expect(ServerAddress.tryParse('host', 'k'), isNull);
  });
}
