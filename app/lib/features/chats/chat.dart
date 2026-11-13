import 'chat_kind.dart';
import 'chat_message.dart';
import 'chat_origin.dart';
import 'chat_role.dart';

class Chat {
  const Chat({
    required this.id,
    required this.kind,
    required this.title,
    required this.about,
    required this.origin,
    required this.peerId,
    required this.role,
    required this.muted,
    required this.request,
    required this.seq,
    required this.rev,
    required this.readSeq,
    required this.unread,
    required this.memberCount,
    required this.last,
    required this.archivedAt,
    required this.changedAt,
  });

  factory Chat.fromJson(Map<String, Object?> json) {
    final origin = json['origin'] as Map<String, Object?>?;
    final last = json['last'] as Map<String, Object?>?;
    final archived = json['archivedAt'] as String?;
    return Chat(
      id: json['id']! as String,
      kind: ChatKind.parse(json['kind']! as String),
      title: json['title'] as String? ?? '',
      about: json['about'] as String? ?? '',
      origin: origin == null ? null : ChatOrigin.fromJson(origin),
      peerId: json['peerId'] as String?,
      role: ChatRole.parse(json['role']! as String),
      muted: json['muted'] as bool? ?? false,
      request: json['request'] as bool? ?? false,
      seq: (json['seq']! as num).toInt(),
      rev: (json['rev']! as num).toInt(),
      readSeq: (json['readSeq']! as num).toInt(),
      unread: (json['unread']! as num).toInt(),
      memberCount: (json['memberCount']! as num).toInt(),
      last: last == null ? null : ChatMessage.fromJson(last),
      archivedAt: archived == null ? null : DateTime.parse(archived),
      changedAt: DateTime.parse(json['changedAt']! as String),
    );
  }

  final String id;
  final ChatKind kind;
  final String title;
  final String about;
  final ChatOrigin? origin;
  final String? peerId;
  final ChatRole role;
  final bool muted;
  final bool request;
  final int seq;
  final int rev;
  final int readSeq;
  final int unread;
  final int memberCount;
  final ChatMessage? last;
  final DateTime? archivedAt;
  final DateTime changedAt;

  bool get archived => archivedAt != null;

  bool get writable => !archived && (kind != ChatKind.channel || role.manages);

  Chat read(int upTo) {
    final seen = upTo > readSeq ? (upTo > seq ? seq : upTo) : readSeq;
    return Chat.fromJson({
      ...toJson(),
      'readSeq': seen,
      'unread': seq - seen < 0 ? 0 : seq - seen,
    });
  }

  Map<String, Object?> toJson() => {
    'id': id,
    'kind': kind.name,
    'title': title,
    'about': about,
    'origin': origin?.toJson(),
    'peerId': peerId,
    'role': role.name,
    'muted': muted,
    'request': request,
    'seq': seq,
    'rev': rev,
    'readSeq': readSeq,
    'unread': unread,
    'memberCount': memberCount,
    'last': last?.toJson(),
    'archivedAt': archivedAt?.toIso8601String(),
    'changedAt': changedAt.toIso8601String(),
  };
}
