import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'modules_section.dart';
import 'settings_section.dart';

class TenantSettingsScreen extends StatelessWidget {
  const TenantSettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.adminSettingsTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: const [
                SettingsSection(),
                SizedBox(height: Spacing.xl),
                ModulesSection(),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
