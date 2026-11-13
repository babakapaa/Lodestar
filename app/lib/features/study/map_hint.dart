class MapHint {
  const MapHint._();

  static const _minWords = 4;

  static bool vague(String text) {
    final words = text.trim().split(RegExp(r'\s+')).where((w) => w.isNotEmpty);
    return words.isNotEmpty && words.length < _minWords;
  }
}
