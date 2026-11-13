import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'leave_request.dart';
import 'team_action.dart';
import 'team_providers.dart';

class LeaveRequestTile extends ConsumerWidget {
  const LeaveRequestTile({
    required this.request,
    required this.onDecided,
    super.key,
  });

  final LeaveRequest request;
  final VoidCallback onDecided;

  Future<void> _decide(
    BuildContext context,
    WidgetRef ref, {
    required bool approve,
  }) async {
    await TeamAction.run(
      context,
      () => ref.read(teamApiProvider).decide(request.id, approve: approve),
    );
    onDecided();
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: const Icon(Icons.logout),
      title: Text(
        l10n.teamLeaveRequestOf(request.student.label, request.teamName),
      ),
      subtitle: request.reason.isEmpty ? null : Text(request.reason),
      trailing: Wrap(
        children: [
          IconButton(
            key: ValueKey('leave-approve-${request.id}'),
            tooltip: l10n.teamLeaveApprove,
            icon: const Icon(Icons.check),
            onPressed: () => _decide(context, ref, approve: true),
          ),
          IconButton(
            tooltip: l10n.teamLeaveKeep,
            icon: const Icon(Icons.close),
            onPressed: () => _decide(context, ref, approve: false),
          ),
        ],
      ),
    );
  }
}
