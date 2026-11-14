import 'package:flutter/foundation.dart';

@immutable
class NotesTabs {
  const NotesTabs({this.open = const [], this.active});

  static const graph = '#graph';

  final List<String> open;
  final String? active;
}
