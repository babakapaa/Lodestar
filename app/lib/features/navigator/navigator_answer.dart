import 'navigator_message.dart';

class NavigatorAnswer {
  const NavigatorAnswer({required this.question, required this.reply});

  factory NavigatorAnswer.fromJson(Map<String, Object?> json) =>
      NavigatorAnswer(
        question: NavigatorMessage.fromJson(
          json['question']! as Map<String, Object?>,
        ),
        reply: NavigatorMessage.fromJson(
          json['reply']! as Map<String, Object?>,
        ),
      );

  final NavigatorMessage question;
  final NavigatorMessage reply;
}
