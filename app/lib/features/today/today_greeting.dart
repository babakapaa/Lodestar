import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/generated/app_localizations.dart';
import '../../l10n/l10n_context.dart';
import '../profile/profile_providers.dart';
import '../study/study_day.dart';
import '../study/study_labels_context.dart';

class TodayGreeting extends ConsumerWidget {
  const TodayGreeting({required this.date, super.key});

  final String date;

  static String hello(AppLocalizations l10n, int hour) {
    if (hour >= 5 && hour < 12) return l10n.todayMorning;
    if (hour >= 12 && hour < 17) return l10n.todayAfternoon;
    if (hour >= 17 && hour < 23) return l10n.todayEvening;
    return l10n.todayNight;
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final textTheme = Theme.of(context).textTheme;
    final name =
        ref.watch(accountProvider).value?.displayName.split(' ').first ?? '';
    final greeting = hello(context.l10n, DateTime.now().hour);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          name.isEmpty ? greeting : '$greeting, $name',
          style: textTheme.headlineMedium,
        ),
        const SizedBox(height: Spacing.xs),
        Text(
          context.studyLabels.day(StudyDay.parse(date)),
          style: textTheme.bodyLarge!.copyWith(
            color: context.tokens.palette.muted,
          ),
        ),
      ],
    );
  }
}
