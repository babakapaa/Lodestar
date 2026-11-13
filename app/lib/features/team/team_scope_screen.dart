import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'team_providers.dart';
import 'team_scope_header.dart';
import 'team_scope_lead.dart';
import 'team_scope_student.dart';

class TeamScopeScreen extends ConsumerWidget {
  const TeamScopeScreen({required this.kind, required this.scopeId, super.key});

  final String kind;
  final String scopeId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final key = (kind, scopeId);
    final view = ref.watch(teamScopeProvider(key));
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.teamsTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: view,
          onRefresh: () => ref.refresh(teamScopeProvider(key).future),
          builder: (view) => [
            TeamScopeHeader(scope: view.scope),
            const SizedBox(height: Spacing.xl),
            if (view.leads)
              TeamScopeLead(view: view)
            else
              TeamScopeStudent(view: view),
          ],
        ),
      ),
    );
  }
}
