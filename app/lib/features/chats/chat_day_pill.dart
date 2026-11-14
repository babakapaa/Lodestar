import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';

class ChatDayPill extends StatelessWidget {
  const ChatDayPill({required this.day, super.key});

  final DateTime day;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = context.palette;
    final now = DateTime.now();
    final locale = Localizations.localeOf(context).toLanguageTag();
    final label = DateUtils.isSameDay(day, now)
        ? l10n.chatsToday
        : DateUtils.isSameDay(day, now.subtract(const Duration(days: 1)))
        ? l10n.chatsYesterday
        : day.year == now.year
        ? DateFormat.MMMMd(locale).format(day)
        : DateFormat.yMMMMd(locale).format(day);
    return Center(
      child: Container(
        margin: const EdgeInsets.symmetric(vertical: Spacing.s),
        padding: const EdgeInsets.symmetric(
          horizontal: Spacing.m,
          vertical: Spacing.xs,
        ),
        decoration: BoxDecoration(
          color: palette.surface,
          borderRadius: BorderRadius.circular(999),
          border: Border.all(color: palette.border),
        ),
        child: Text(
          label,
          style: Theme.of(context).textTheme.labelSmall!
              .copyWith(color: palette.muted),
        ),
      ),
    );
  }
}
