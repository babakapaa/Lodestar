class CanvasEdge {
  const CanvasEdge({required this.id, required this.from, required this.to});

  factory CanvasEdge.fromJson(Map<String, Object?> json) => CanvasEdge(
    id: json['id']! as String,
    from: json['fromNode']! as String,
    to: json['toNode']! as String,
  );

  final String id;
  final String from;
  final String to;

  Map<String, Object?> toJson() => {
    'id': id,
    'fromNode': from,
    'fromSide': 'right',
    'toNode': to,
    'toSide': 'left',
  };
}
