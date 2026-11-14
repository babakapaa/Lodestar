import 'package:material_ui/material_ui.dart';

abstract final class BubbleShape {
  static const _round = Radius.circular(16);
  static const _tail = Radius.circular(4);

  static BorderRadius of({required bool mine, required bool grouped}) =>
      BorderRadiusDirectional.only(
        topStart: !mine && grouped ? _tail : _round,
        topEnd: mine && grouped ? _tail : _round,
        bottomStart: mine ? _round : _tail,
        bottomEnd: mine ? _tail : _round,
      ).resolve(TextDirection.ltr);
}
