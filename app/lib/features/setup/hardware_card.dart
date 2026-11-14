import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/section.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'hardware_check.dart';
import 'hardware_tile.dart';
import 'setup_providers.dart';

class HardwareCard extends ConsumerStatefulWidget {
  const HardwareCard({super.key});

  @override
  ConsumerState<HardwareCard> createState() => _HardwareCardState();
}

class _HardwareCardState extends ConsumerState<HardwareCard> {
  bool _busy = false;

  Future<void> _check() async {
    setState(() => _busy = true);
    await TeamAction.run(context, () async {
      await ref.read(setupApiProvider).checkHardware();
      ref
        ..invalidate(setupHardwareProvider)
        ..invalidate(setupOverviewProvider);
    });
    if (mounted) {
      setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final report = ref.watch(setupHardwareProvider).value;
    final checkedAt = report?.checkedAt;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    return Section(
      title: l10n.setupHardware,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            checkedAt == null
                ? l10n.setupHardwareNever
                : l10n.setupHardwareCheckedAt(
                    DateFormat.yMd(l10n.localeName).add_Hm().format(checkedAt),
                  ),
            style: muted,
          ),
          for (final check in report?.checks ?? const <HardwareCheck>[])
            HardwareTile(check: check),
          Align(
            alignment: Alignment.centerLeft,
            child: TextButton.icon(
              key: const ValueKey('setup-hardware-check'),
              onPressed: _busy ? null : _check,
              icon: const Icon(Icons.refresh),
              label: Text(l10n.setupHardwareRun),
            ),
          ),
        ],
      ),
    );
  }
}
