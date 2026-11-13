import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/lodestar_card.dart';
import 'run_icon.dart';
import 'sandbox_labels_context.dart';
import 'sandbox_run.dart';

class RunTile extends StatelessWidget {
  const RunTile({required this.run, super.key});

  final SandboxRun run;

  @override
  Widget build(BuildContext context) {
    return LodestarCard(
      child: ListTile(
        contentPadding: EdgeInsets.zero,
        leading: Icon(RunIcon.of(run.kind)),
        title: Text(run.title),
        subtitle: Text(context.sandboxLabels.summary(run)),
        trailing: const Icon(Icons.chevron_right),
        onTap: () => context.push(AppRoutes.sandboxRun(run.id)),
      ),
    );
  }
}
