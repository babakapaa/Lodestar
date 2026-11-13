import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../l10n/l10n_context.dart';
import 'team_profile_form.dart';
import 'team_providers.dart';

class TeamProfileScreen extends ConsumerWidget {
  const TeamProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final profile = ref.watch(teamProfileProvider);
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.teamProfileTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: profile,
          onRefresh: () => ref.refresh(teamProfileProvider.future),
          builder: (profile) => [TeamProfileForm(profile: profile)],
        ),
      ),
    );
  }
}
