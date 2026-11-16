import 'package:flutter/foundation.dart';

enum MdKind { heading, paragraph, item, quote, rule, code, math }

@immutable
class MdBlock {
  const MdBlock(
    this.kind,
    this.text, {
    this.level = 0,
    this.marker = '',
    this.done,
  });

  final MdKind kind;
  final String text;
  final int level;
  final String marker;
  final bool? done;
}
