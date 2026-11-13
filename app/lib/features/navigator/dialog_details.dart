import 'brief.dart';
import 'navigator_dialog.dart';
import 'navigator_message.dart';

class DialogDetails {
  const DialogDetails({
    required this.dialog,
    required this.brief,
    required this.messages,
  });

  factory DialogDetails.fromJson(Map<String, Object?> json) => DialogDetails(
    dialog: NavigatorDialog.fromJson(json),
    brief: Brief.maybe(json['brief']),
    messages: [
      for (final item in json['messages']! as List<Object?>)
        NavigatorMessage.fromJson(item! as Map<String, Object?>),
    ],
  );

  final NavigatorDialog dialog;
  final Brief? brief;
  final List<NavigatorMessage> messages;
}
