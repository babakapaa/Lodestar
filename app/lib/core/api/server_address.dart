class ServerAddress {
  const ServerAddress({required this.server, required this.tenant});

  static final _tenantPattern = RegExp(r'^[a-z0-9][a-z0-9-]{1,39}$');

  final Uri server;
  final String tenant;

  static ServerAddress? tryParse(String server, String tenant) {
    var text = server.trim();
    if (text.isEmpty) {
      return null;
    }
    if (!text.contains('://')) {
      text = 'https://$text';
    }
    final uri = Uri.tryParse(text);
    if (uri == null ||
        uri.host.isEmpty ||
        (uri.scheme != 'https' && uri.scheme != 'http')) {
      return null;
    }
    final slug = tenant.trim().toLowerCase();
    if (!_tenantPattern.hasMatch(slug)) {
      return null;
    }
    var path = uri.path;
    while (path.endsWith('/')) {
      path = path.substring(0, path.length - 1);
    }
    return ServerAddress(
      server: Uri(
        scheme: uri.scheme,
        host: uri.host,
        port: uri.hasPort ? uri.port : null,
        path: path,
      ),
      tenant: slug,
    );
  }

  Uri resolve(String path, [Map<String, String>? query]) {
    return server.replace(path: '${server.path}$path', queryParameters: query);
  }

  @override
  bool operator ==(Object other) =>
      other is ServerAddress &&
      other.server == server &&
      other.tenant == tenant;

  @override
  int get hashCode => Object.hash(server, tenant);
}
