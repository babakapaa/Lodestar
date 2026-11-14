import 'package:flutter/foundation.dart';

@immutable
class NoteHeading {
  const NoteHeading(this.level, this.title, this.offset);

  static final _pattern = RegExp(r'^(#{1,6})\s+(.+)$', multiLine: true);

  static List<NoteHeading> parse(String body) {
    final headings = <NoteHeading>[];
    var fenced = false;
    var offset = 0;
    for (final line in body.split('\n')) {
      if (line.trimLeft().startsWith('```')) {
        fenced = !fenced;
      } else if (!fenced) {
        final match = _pattern.firstMatch(line);
        if (match != null) {
          headings.add(
            NoteHeading(match.group(1)!.length, match.group(2)!.trim(), offset),
          );
        }
      }
      offset += line.length + 1;
    }
    return headings;
  }

  final int level;
  final String title;
  final int offset;
}
