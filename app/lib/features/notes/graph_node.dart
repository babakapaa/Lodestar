import 'dart:ui';

class GraphNode {
  GraphNode(this.id, this.label, this.position);

  final String id;
  final String label;
  Offset position;
  int degree = 0;
}
