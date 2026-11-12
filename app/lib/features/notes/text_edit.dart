class TextEdit {
  const TextEdit(this.index, this.remove, this.insert);

  final int index;
  final int remove;
  final String insert;

  static TextEdit? between(String before, String after) {
    if (before == after) {
      return null;
    }
    final limit = before.length < after.length ? before.length : after.length;
    var start = 0;
    while (start < limit &&
        before.codeUnitAt(start) == after.codeUnitAt(start)) {
      start++;
    }
    var end = 0;
    while (end < limit - start &&
        before.codeUnitAt(before.length - 1 - end) ==
            after.codeUnitAt(after.length - 1 - end)) {
      end++;
    }
    if (start > 0 && _isHighSurrogate(before.codeUnitAt(start - 1))) {
      start--;
    }
    if (end > 0 && _isLowSurrogate(before.codeUnitAt(before.length - end))) {
      end--;
    }
    return TextEdit(
      start,
      before.length - start - end,
      after.substring(start, after.length - end),
    );
  }

  int shift(int offset) {
    if (offset <= index) {
      return offset;
    }
    if (offset >= index + remove) {
      return offset - remove + insert.length;
    }
    return index + insert.length;
  }

  static bool _isHighSurrogate(int unit) => unit >= 0xD800 && unit <= 0xDBFF;

  static bool _isLowSurrogate(int unit) => unit >= 0xDC00 && unit <= 0xDFFF;
}
