import 'dart:convert';
import 'dart:math';
import 'dart:ui';

import 'canvas_card.dart';
import 'canvas_edge.dart';

class CanvasBoard {
  CanvasBoard(this.cards, this.edges);

  factory CanvasBoard.parse(String text) {
    if (text.trim().isEmpty) {
      return CanvasBoard([], []);
    }
    try {
      final json = jsonDecode(text) as Map<String, Object?>;
      final cards = [
        for (final item in (json['nodes'] as List<Object?>?) ?? const [])
          if (item is Map<String, Object?> &&
              (item['type'] == 'text' || item['type'] == 'file'))
            CanvasCard.fromJson(item),
      ];
      final ids = {for (final card in cards) card.id};
      final edges = [
        for (final item in (json['edges'] as List<Object?>?) ?? const [])
          if (item is Map<String, Object?>) CanvasEdge.fromJson(item),
      ].where((edge) => ids.contains(edge.from) && ids.contains(edge.to));
      return CanvasBoard(cards, edges.toList());
    } on FormatException {
      return CanvasBoard([], []);
    } on TypeError {
      return CanvasBoard([], []);
    }
  }

  static final _random = Random();

  static String newId() =>
      List.generate(16, (_) => _random.nextInt(16).toRadixString(16)).join();

  final List<CanvasCard> cards;
  final List<CanvasEdge> edges;

  CanvasCard? card(String id) {
    for (final card in cards) {
      if (card.id == id) {
        return card;
      }
    }
    return null;
  }

  CanvasCard? cardAt(Offset point) {
    for (final card in cards.reversed) {
      if (card.rect.contains(point)) {
        return card;
      }
    }
    return null;
  }

  void remove(String id) {
    cards.removeWhere((card) => card.id == id);
    edges.removeWhere(
      (edge) => edge.id == id || edge.from == id || edge.to == id,
    );
  }

  void connect(String from, String to) {
    if (from == to ||
        edges.any(
          (edge) =>
              (edge.from == from && edge.to == to) ||
              (edge.from == to && edge.to == from),
        )) {
      return;
    }
    edges.add(CanvasEdge(id: newId(), from: from, to: to));
  }

  Rect? get bounds {
    if (cards.isEmpty) {
      return null;
    }
    return cards
        .map((card) => card.rect)
        .reduce((a, b) => a.expandToInclude(b));
  }

  String encode() => const JsonEncoder.withIndent('\t').convert({
    'nodes': [for (final card in cards) card.toJson()],
    'edges': [for (final edge in edges) edge.toJson()],
  });
}
