import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/auth/auth_providers.dart';
import '../../core/auth/auth_state.dart';
import '../../core/navigation/app_routes.dart';
import '../../core/navigation/entry_page.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'lti_launch.dart';

class LtiLaunchScreen extends ConsumerStatefulWidget {
  const LtiLaunchScreen({required this.launch, super.key});

  final LtiLaunch launch;

  @override
  ConsumerState<LtiLaunchScreen> createState() => _LtiLaunchScreenState();
}

class _LtiLaunchScreenState extends ConsumerState<LtiLaunchScreen> {
  String? _error;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    final launch = widget.launch;
    final auth = ref.read(authProvider);
    if (launch.valid &&
        (auth is Disconnected || auth.address == launch.address)) {
      unawaited(Future.microtask(_enter));
    }
  }

  Future<void> _enter() async {
    final launch = widget.launch;
    final l10n = context.l10n;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref
          .read(authProvider.notifier)
          .claim(launch.address!, launch.ticket!);
      if (mounted) {
        context.go(launch.target);
      }
    } on Exception catch (error) {
      if (mounted) {
        setState(() {
          _busy = false;
          _error = ApiErrorText.of(l10n, error);
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final launch = widget.launch;
    final address = launch.address;
    final error = _error;
    if (!launch.valid) {
      return EntryPage(
        title: l10n.ltiTitle,
        subtitle: l10n.ltiFailed,
        child: LodestarButton(
          label: l10n.ltiHome,
          variant: LodestarButtonVariant.secondary,
          onPressed: () => context.go(AppRoutes.todayPath),
        ),
      );
    }
    return EntryPage(
      title: l10n.ltiTitle,
      subtitle: l10n.loginBody(address!.tenant, address.server.host),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (_busy)
            const Center(child: CircularProgressIndicator())
          else ...[
            Text(l10n.ltiConfirm),
            const SizedBox(height: Spacing.m),
            LodestarButton(
              key: const ValueKey('lti-enter'),
              label: l10n.ltiEnter,
              icon: Icons.login,
              onPressed: _enter,
            ),
            const SizedBox(height: Spacing.s),
            TextButton(
              onPressed: () => context.go(AppRoutes.todayPath),
              child: Text(l10n.actionCancel),
            ),
          ],
          if (error != null) ...[
            const SizedBox(height: Spacing.m),
            FormError(message: error),
          ],
        ],
      ),
    );
  }
}
