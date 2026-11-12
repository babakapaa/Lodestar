import 'package:material_ui/material_ui.dart';

abstract final class RoleColor {
  static const choices = [
    '',
    '#C0392B',
    '#D35400',
    '#B7950B',
    '#1E8449',
    '#0B786E',
    '#2563EB',
    '#7D3C98',
  ];

  static Color? parse(String value) {
    if (value.length != 7 || !value.startsWith('#')) {
      return null;
    }
    final rgb = int.tryParse(value.substring(1), radix: 16);
    return rgb == null ? null : Color(0xFF000000 | rgb);
  }
}
