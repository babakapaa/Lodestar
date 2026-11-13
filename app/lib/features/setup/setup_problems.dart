import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../l10n/l10n_context.dart';
import 'setup_labels.dart';
import 'setup_problem.dart';

class SetupProblems extends StatelessWidget {
  const SetupProblems({required this.problems, super.key});

  final List<SetupProblem> problems;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = SetupLabels(l10n);
    return Section(
      title: l10n.setupProblemsTitle,
      child: problems.isEmpty
          ? Text(l10n.setupReady)
          : Column(
              children: [
                for (final p in problems)
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.error_outline),
                    title: Text(labels.problem(p)),
                    trailing: TextButton(
                      onPressed: () =>
                          context.push(AppRoutes.setupSection(p.section)),
                      child: Text(l10n.setupFix),
                    ),
                  ),
              ],
            ),
    );
  }
}
