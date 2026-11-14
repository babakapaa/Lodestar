import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/api/connection_store.dart';
import '../../core/api/server_address.dart';
import '../../core/auth/auth_providers.dart';
import '../../core/navigation/entry_page.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';

class ConnectScreen extends ConsumerStatefulWidget {
  const ConnectScreen({super.key});

  @override
  ConsumerState<ConnectScreen> createState() => _ConnectScreenState();
}

class _ConnectScreenState extends ConsumerState<ConnectScreen> {
  final _server = TextEditingController(text: ConnectionStore.suggestedServer);
  final _tenant = TextEditingController(text: ConnectionStore.presetTenant);
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    _server.dispose();
    _tenant.dispose();
    super.dispose();
  }

  Future<void> _connect() async {
    final l10n = context.l10n;
    final address = ServerAddress.tryParse(_server.text, _tenant.text);
    if (address == null) {
      setState(() => _error = l10n.connectInvalid);
      return;
    }
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref.read(authProvider.notifier).connect(address);
    } on ApiException catch (error) {
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
    final error = _error;
    return EntryPage(
      title: l10n.connectTitle,
      subtitle: l10n.connectBody,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          TextField(
            controller: _server,
            enabled: !_busy,
            keyboardType: TextInputType.url,
            autocorrect: false,
            textInputAction: TextInputAction.next,
            decoration: InputDecoration(
              labelText: l10n.connectServer,
              hintText: 'campus.example.edu',
            ),
          ),
          const SizedBox(height: Spacing.m),
          TextField(
            controller: _tenant,
            enabled: !_busy,
            autocorrect: false,
            textInputAction: TextInputAction.done,
            onSubmitted: (_) => _connect(),
            decoration: InputDecoration(
              labelText: l10n.connectTenant,
              helperText: l10n.connectTenantHint,
            ),
          ),
          if (error != null) FormError(message: error),
          const SizedBox(height: Spacing.l),
          LodestarButton(
            label: l10n.connectAction,
            onPressed: _busy ? null : _connect,
          ),
        ],
      ),
    );
  }
}
