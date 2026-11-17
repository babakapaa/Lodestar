import 'package:flutter/foundation.dart';

enum MdKind { heading, paragraph, item, quote, rule, code, math, table }

@immutable
class MdBlock {
  const MdBlock(
    this.kind,
    this.text, {
    this.level = 0,
    this.marker = '',
    this.done,
    this.rows = const [],
  });

  final MdKind kind;
  final String text;
  final int level;
  final String marker;
  final bool? done;
  final List<List<String>> rows;
}
