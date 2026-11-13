import 'question_kind.dart';

class DefenseQuestion {
  const DefenseQuestion({
    required this.kind,
    required this.text,
    this.mark,
    this.note = '',
  });

  factory DefenseQuestion.fromJson(Map<String, Object?> json) =>
      DefenseQuestion(
        kind: QuestionKind.parse(json['kind']),
        text: json['text'] as String? ?? '',
        mark: (json['mark'] as num?)?.toInt(),
        note: json['note'] as String? ?? '',
      );

  static const maxMark = 3;

  final QuestionKind kind;
  final String text;
  final int? mark;
  final String note;

  DefenseQuestion copyWith({String? text, int? mark, String? note}) =>
      DefenseQuestion(
        kind: kind,
        text: text ?? this.text,
        mark: mark ?? this.mark,
        note: note ?? this.note,
      );

  Map<String, Object?> toJson() => {
    'kind': kind.name,
    'text': text,
    'mark': mark,
    'note': note,
  };
}
