import 'package:material_ui/material_ui.dart';

class ChatAvatar extends StatelessWidget {
  const ChatAvatar({required this.title, this.size = 44, super.key});

  static const _hues = [4.0, 28.0, 45.0, 140.0, 175.0, 200.0, 230.0, 275.0];

  final String title;
  final double size;

  String get _initials {
    final words = title
        .split(RegExp(r'\s+'))
        .where((word) => word.isNotEmpty)
        .toList();
    if (words.isEmpty) {
      return '?';
    }
    final first = words.first.characters.first;
    final second = words.length > 1 ? words[1].characters.first : '';
    return (first + second).toUpperCase();
  }

  @override
  Widget build(BuildContext context) {
    final hue = _hues[title.hashCode.abs() % _hues.length];
    final dark = Theme.of(context).brightness == Brightness.dark;
    final color = HSLColor.fromAHSL(1, hue, 0.55, dark ? 0.42 : 0.5).toColor();
    return Container(
      width: size,
      height: size,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            HSLColor.fromColor(color).withLightness(dark ? 0.5 : 0.6).toColor(),
            color,
          ],
        ),
      ),
      child: Text(
        _initials,
        style: TextStyle(
          color: Colors.white,
          fontWeight: FontWeight.w600,
          fontSize: size * 0.36,
        ),
      ),
    );
  }
}
