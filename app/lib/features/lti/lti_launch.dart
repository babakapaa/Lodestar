import '../../core/api/server_address.dart';
import '../../core/navigation/app_routes.dart';

class LtiLaunch {
  const LtiLaunch({this.address, this.ticket, this.next, this.failed = false});

  factory LtiLaunch.parse(Map<String, String> query) {
    final server = query['server'];
    final tenant = query['tenant'];
    final next = query['next'];
    return LtiLaunch(
      address: server == null || tenant == null
          ? null
          : ServerAddress.tryParse(server, tenant),
      ticket: query['ticket'],
      next: next != null && next.startsWith('/') && !next.startsWith('//')
          ? next
          : null,
      failed: query.containsKey('error'),
    );
  }

  final ServerAddress? address;
  final String? ticket;
  final String? next;
  final bool failed;

  bool get valid => !failed && address != null && (ticket?.isNotEmpty ?? false);

  String get target => next ?? AppRoutes.todayPath;
}
