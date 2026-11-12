import 'dart:math';
import 'dart:ui';

import 'graph_node.dart';
import 'note.dart';

class NotesGraph {
  NotesGraph._(this.nodes, this.edges);

  factory NotesGraph.of(List<Note> notes) {
    final live = [
      for (final note in notes)
        if (!note.deleted) note,
    ];
    final byKey = <String, int>{};
    for (var i = 0; i < live.length; i++) {
      byKey.putIfAbsent(live[i].meta.path, () => i);
      byKey.putIfAbsent(live[i].meta.name, () => i);
    }
    final nodes = [
      for (var i = 0; i < live.length; i++)
        GraphNode(
          live[i].id,
          live[i].meta.name,
          Offset(
            cos(2 * pi * i / max(live.length, 1)) * 200,
            sin(2 * pi * i / max(live.length, 1)) * 200,
          ),
        ),
    ];
    final edges = <(int, int)>{};
    for (var i = 0; i < live.length; i++) {
      for (final link in live[i].meta.links) {
        final j = byKey[link];
        if (j != null && j != i) {
          edges.add(i < j ? (i, j) : (j, i));
        }
      }
    }
    for (final (a, b) in edges) {
      nodes[a].degree++;
      nodes[b].degree++;
    }
    return NotesGraph._(nodes, edges.toList())..settle();
  }

  static const _spring = 120.0;

  final List<GraphNode> nodes;
  final List<(int, int)> edges;

  void settle({int steps = 300}) {
    final n = nodes.length;
    if (n < 2) {
      return;
    }
    for (var step = 0; step < steps; step++) {
      final cooling = 1 - step / steps;
      final forces = List.filled(n, Offset.zero);
      for (var i = 0; i < n; i++) {
        for (var j = i + 1; j < n; j++) {
          var delta = nodes[i].position - nodes[j].position;
          if (delta.distance < 0.01) {
            delta = Offset(0.01 * (i + 1), 0.01 * (j + 1));
          }
          final push =
              delta / delta.distance * (_spring * _spring) / delta.distance;
          forces[i] += push;
          forces[j] -= push;
        }
      }
      for (final (a, b) in edges) {
        final delta = nodes[a].position - nodes[b].position;
        final pull =
            delta * delta.distance / _spring / max(delta.distance, 0.01);
        forces[a] -= pull;
        forces[b] += pull;
      }
      for (var i = 0; i < n; i++) {
        final pullToCenter = nodes[i].position * 0.02;
        var move = forces[i] - pullToCenter;
        final limit = 20 * cooling + 0.5;
        if (move.distance > limit) {
          move = move / move.distance * limit;
        }
        nodes[i].position += move;
      }
    }
  }

  Rect get bounds {
    if (nodes.isEmpty) {
      return Rect.zero;
    }
    var rect = Rect.fromCircle(center: nodes.first.position, radius: 1);
    for (final node in nodes) {
      rect = rect.expandToInclude(
        Rect.fromCircle(center: node.position, radius: 1),
      );
    }
    return rect;
  }

  GraphNode? hit(Offset point, double radius) {
    for (final node in nodes) {
      if ((node.position - point).distance <= radius) {
        return node;
      }
    }
    return null;
  }
}
