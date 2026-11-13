import 'call.dart';

class CallTicket {
  const CallTicket({
    required this.call,
    required this.url,
    required this.token,
    required this.screenShare,
  });

  factory CallTicket.fromJson(Map<String, Object?> json) => CallTicket(
    call: Call.fromJson(json['call']! as Map<String, Object?>),
    url: json['url']! as String,
    token: json['token']! as String,
    screenShare: json['screenShare']! as bool,
  );

  final Call call;
  final String url;
  final String token;
  final bool screenShare;
}
