import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/notice_banner.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'leave_status.dart';
import 'team_action.dart';
import 'team_details.dart';
import 'team_providers.dart';
import 'team_refresh.dart';
import 'team_text_dialog.dart';

class TeamLeaveSection extends ConsumerWidget {
  const TeamLeaveSection({required this.details, super.key});

  final TeamDetails details;

  Future<void> _leave(BuildContext context, WidgetRef ref) async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    final free = details.canLeaveFreely;
    final reason = await TeamTextDialog.show(
      context,
      title: free ? l10n.teamLeave : l10n.teamLeaveAsk,
      label: l10n.teamLeaveReason,
      hint: free ? l10n.teamLeaveFreeHint : l10n.teamLeaveAskHint,
      action: free ? l10n.teamLeave : l10n.teamLeaveSend,
      required: !free,
      lines: 2,
    );
    if (reason == null || !context.mounted) {
      return;
    }
    var left = false;
    await TeamAction.run(context, () async {
      left =
          (await ref.read(teamApiProvider).leave(details.team.id, reason)).left;
    });
    ref.refreshTeam(details);
    if (left && router.canPop()) {
      router.pop();
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final request = details.request;
    if (details.team.archived) {
      return const SizedBox.shrink();
    }
    if (request != null && request.status == LeaveStatus.pending) {
      return NoticeBanner(icon: Icons.logout, message: l10n.teamLeavePending);
    }
    return Section(
      title: l10n.teamLeave,
      child: Align(
        alignment: Alignment.centerLeft,
        child: LodestarButton(
          key: const ValueKey('team-leave'),
          label: details.canLeaveFreely ? l10n.teamLeave : l10n.teamLeaveAsk,
          icon: Icons.logout,
          variant: LodestarButtonVariant.secondary,
          onPressed: () => _leave(context, ref),
        ),
      ),
    );
  }
}
