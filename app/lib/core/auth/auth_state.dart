import '../api/server_address.dart';

sealed class AuthState {
  const AuthState();

  ServerAddress? get address;
}

final class Disconnected extends AuthState {
  const Disconnected();

  @override
  ServerAddress? get address => null;
}

final class SignedOut extends AuthState {
  const SignedOut(this.address);

  @override
  final ServerAddress address;
}

final class SignedIn extends AuthState {
  const SignedIn(this.address);

  @override
  final ServerAddress address;
}
