import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../people/people_search.dart';

class PeopleRolesScreen extends StatelessWidget {
  const PeopleRolesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.peopleRolesTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: [
                PeopleSearch(
                  onSelected: (person) =>
                      context.push(AppRoutes.personRoles(person.id)),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
