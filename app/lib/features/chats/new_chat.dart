import 'chat_kind.dart';
import 'chat_origin.dart';

class NewChat {
  const NewChat({
    required this.kind,
    this.title = '',
    this.about = '',
    this.memberIds = const [],
    this.origin,
  });

  final ChatKind kind;
  final String title;
  final String about;
  final List<String> memberIds;
  final ChatOrigin? origin;

  Map<String, Object?> toJson() => {
    'kind': kind.name,
    'title': title,
    'about': about,
    'memberIds': memberIds,
    'origin': ?origin?.toJson(),
  };
}
