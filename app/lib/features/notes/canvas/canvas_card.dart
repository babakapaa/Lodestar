import 'dart:ui';

class CanvasCard {
  CanvasCard({required this.id, required this.rect, this.text = '', this.file});

  factory CanvasCard.fromJson(Map<String, Object?> json) => CanvasCard(
    id: json['id']! as String,
    rect: Rect.fromLTWH(
      (json['x'] as num? ?? 0).toDouble(),
      (json['y'] as num? ?? 0).toDouble(),
      (json['width'] as num? ?? 260).toDouble(),
      (json['height'] as num? ?? 140).toDouble(),
    ),
    text: json['text'] as String? ?? '',
    file: json['type'] == 'file' ? json['file'] as String? : null,
  );

  static const minSize = Size(120, 60);

  final String id;
  Rect rect;
  String text;
  final String? file;

  bool get isNote => file != null;

  String? get notePath {
    final path = file;
    if (path == null) {
      return null;
    }
    return path.toLowerCase().endsWith('.md')
        ? path.substring(0, path.length - 3)
        : path;
  }

  Map<String, Object?> toJson() => {
    'id': id,
    'type': isNote ? 'file' : 'text',
    if (isNote) 'file': file else 'text': text,
    'x': rect.left.round(),
    'y': rect.top.round(),
    'width': rect.width.round(),
    'height': rect.height.round(),
  };
}
