import 'package:material_ui/material_ui.dart';

class MentionQuery {
  const MentionQuery(this.start, this.end, this.text);

  static final _pattern = RegExp(r'(^|\s)@([^\s@]*)$');

  final int start;
  final int end;
  final String text;

  static MentionQuery? of(TextEditingValue value) {
    final cursor = value.selection.baseOffset;
    if (!value.selection.isCollapsed || cursor < 0) {
      return null;
    }
    final match = _pattern.firstMatch(value.text.substring(0, cursor));
    if (match == null) {
      return null;
    }
    final query = match.group(2)!;
    return MentionQuery(cursor - query.length - 1, cursor, query);
  }

  TextEditingValue insert(TextEditingValue value, String name) {
    final mention = '@$name ';
    return TextEditingValue(
      text: value.text.replaceRange(start, end, mention),
      selection: TextSelection.collapsed(offset: start + mention.length),
    );
  }
}
