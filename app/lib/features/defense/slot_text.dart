import 'package:intl/intl.dart';

import 'defense_slot.dart';

abstract final class SlotText {
  static String of(DefenseSlot slot, String locale) {
    final start = slot.startsAt.toLocal();
    final when =
        '${DateFormat.MMMMEEEEd(locale).format(start)}, '
        '${DateFormat.Hm(locale).format(start)}–'
        '${DateFormat.Hm(locale).format(slot.endsAt.toLocal())}';
    return slot.place.isEmpty ? when : '$when, ${slot.place}';
  }
}
