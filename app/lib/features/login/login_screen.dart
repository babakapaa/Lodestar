import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/auth/auth_providers.dart';
import '../../core/auth/identity_provider.dart';
import '../../core/navigation/entry_page.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'identity_providers_provider.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  String? _error;
  bool _busy = false;

  Future<void> _signIn(IdentityProvider provider) async {
    final l10n = context.l10n;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref.read(authProvider.notifier).signIn(provider);
    } on Exception catch (error) {
      if (mounted) {
        setState(() => _error = ApiErrorText.of(l10n, error));
      }
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final address = ref.watch(authProvider).address;
    final providers = ref.watch(identityProvidersProvider);
    final error = _error;
    return EntryPage(
      title: l10n.loginTitle,
      subtitle: address == null
          ? ''
          : l10n.loginBody(address.tenant, address.server.host),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          ...switch (providers) {
            AsyncData(:final value) when value.isEmpty => [
              Text(l10n.loginNoProviders),
            ],
            AsyncData(:final value) => [
              for (final provider in value)
                Padding(
                  padding: const EdgeInsets.only(bottom: Spacing.s),
                  child: LodestarButton(
                    label: l10n.loginWith(provider.displayName),
                    icon: Icons.login,
                    onPressed: _busy ? null : () => _signIn(provider),
                  ),
                ),
            ],
            AsyncError(:final error) => [
              FormError(message: ApiErrorText.of(l10n, error)),
              const SizedBox(height: Spacing.m),
              LodestarButton(
                label: l10n.actionRetry,
                variant: LodestarButtonVariant.secondary,
                onPressed: () => ref.invalidate(identityProvidersProvider),
              ),
            ],
            _ => [const Center(child: CircularProgressIndicator())],
          },
          if (error != null) FormError(message: error),
          const SizedBox(height: Spacing.l),
          TextButton(
            onPressed: _busy
                ? null
                : () => ref.read(authProvider.notifier).disconnect(),
            child: Text(l10n.loginChangeUniversity),
          ),
        ],
      ),
    );
  }
}
