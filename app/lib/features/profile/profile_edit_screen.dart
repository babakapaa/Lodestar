import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'profile_form.dart';
import 'profile_providers.dart';

class ProfileEditScreen extends ConsumerWidget {
  const ProfileEditScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final profile = ref.watch(myProfileProvider);
    final account = ref.watch(accountProvider);
    final error = profile.error ?? account.error;
    final mine = profile.value;
    final me = account.value;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.profileEdit)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 640),
            child: error != null
                ? Padding(
                    padding: const EdgeInsets.all(Spacing.l),
                    child: FormError(message: ApiErrorText.of(l10n, error)),
                  )
                : mine == null || me == null
                ? const Center(child: CircularProgressIndicator())
                : ProfileForm(profile: mine, fallbackName: me.displayName),
          ),
        ),
      ),
    );
  }
}
