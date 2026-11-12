abstract interface class LiveSocket {
  Stream<String> connect(Uri uri);
}
