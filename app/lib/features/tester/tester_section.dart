import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'tester_people.dart';
import 'tester_providers.dart';

class TesterSection extends ConsumerWidget {
  const TesterSection({required this.hint, this.current, super.key});

  final String hint;
  final String? current;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final info = ref.watch(testerInfoProvider).value;
    if (info == null) {
      return const SizedBox.shrink();
    }
    return Padding(
      padding: const EdgeInsets.only(top: Spacing.xl),
      child: Section(
        title: context.l10n.testerTitle,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(hint),
            const SizedBox(height: Spacing.m),
            TesterPeople(people: info.people, current: current),
          ],
        ),
      ),
    );
  }
}
