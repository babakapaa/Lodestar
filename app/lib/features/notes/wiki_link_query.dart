import 'package:material_ui/material_ui.dart';

abstract final class WikiLinkQuery {
  static final _link = RegExp(r'\[\[([^\[\]\n]+?)\]\]');

  static String? typed(TextEditingValue value) {
    final cursor = value.selection.baseOffset;
    if (!value.selection.isCollapsed || cursor < 2) {
      return null;
    }
    final before = value.text.substring(0, cursor);
    final open = before.lastIndexOf('[[');
    if (open < 0) {
      return null;
    }
    final query = before.substring(open + 2);
    if (query.contains(']') || query.contains('\n') || query.contains('[')) {
      return null;
    }
    return query;
  }

  static TextEditingValue insert(TextEditingValue value, String name) {
    final cursor = value.selection.baseOffset;
    final text = value.text;
    final open = text.substring(0, cursor).lastIndexOf('[[');
    final after = text.substring(cursor);
    final closed = after.startsWith(']]');
    final inserted = '[[$name]]';
    final tail = closed ? after.substring(2) : after;
    final next = text.substring(0, open) + inserted + tail;
    return TextEditingValue(
      text: next,
      selection: TextSelection.collapsed(offset: open + inserted.length),
    );
  }

  static String? at(TextEditingValue value) {
    final cursor = value.selection.baseOffset;
    if (cursor < 0) {
      return null;
    }
    for (final match in _link.allMatches(value.text)) {
      if (match.start <= cursor && cursor <= match.end) {
        final target = match.group(1)!.split('|').first.split('#').first.trim();
        return target.isEmpty ? null : target;
      }
    }
    return null;
  }
}
