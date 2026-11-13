import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'lti_platform.dart';
import 'lti_platform_input.dart';

class LtiPlatformDialog extends StatefulWidget {
  const LtiPlatformDialog({this.platform, super.key});

  final LtiPlatform? platform;

  static Future<LtiPlatformInput?> show(
    BuildContext context, [
    LtiPlatform? platform,
  ]) => showDialog<LtiPlatformInput>(
    context: context,
    builder: (context) => LtiPlatformDialog(platform: platform),
  );

  @override
  State<LtiPlatformDialog> createState() => _LtiPlatformDialogState();
}

class _LtiPlatformDialogState extends State<LtiPlatformDialog> {
  late final LtiPlatformInput? _initial = widget.platform?.input;
  late final _name = TextEditingController(text: _initial?.name ?? '');
  late final _issuer = TextEditingController(text: _initial?.issuer ?? '');
  late final _client = TextEditingController(text: _initial?.clientId ?? '');
  late final _auth = TextEditingController(text: _initial?.authUrl ?? '');
  late final _token = TextEditingController(text: _initial?.tokenUrl ?? '');
  late final _jwks = TextEditingController(text: _initial?.jwksUrl ?? '');
  late final _deployments = TextEditingController(
    text: _initial?.deploymentIds.join(', ') ?? '',
  );
  late bool _matchEmail = _initial?.matchEmail ?? false;

  List<TextEditingController> get _fields => [
    _name,
    _issuer,
    _client,
    _auth,
    _token,
    _jwks,
    _deployments,
  ];

  @override
  void dispose() {
    for (final c in _fields) {
      c.dispose();
    }
    super.dispose();
  }

  void _done() {
    final deployments = [
      for (final id in _deployments.text.split(RegExp(r'[,\s]+')))
        if (id.trim().isNotEmpty) id.trim(),
    ];
    if (_fields.any((c) => c.text.trim().isEmpty) || deployments.isEmpty) {
      return;
    }
    Navigator.of(context).pop(
      LtiPlatformInput(
        name: _name.text.trim(),
        issuer: _issuer.text.trim(),
        clientId: _client.text.trim(),
        authUrl: _auth.text.trim(),
        tokenUrl: _token.text.trim(),
        jwksUrl: _jwks.text.trim(),
        deploymentIds: deployments,
        matchEmail: _matchEmail,
      ),
    );
  }

  Widget _field(String key, TextEditingController controller, String label) =>
      TextField(
        key: ValueKey('lms-$key'),
        controller: controller,
        decoration: InputDecoration(labelText: label),
      );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    return AlertDialog(
      title: Text(widget.platform == null ? l10n.lmsNew : l10n.lmsEdit),
      content: SizedBox(
        width: 420,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              _field('name', _name, l10n.lmsName),
              _field('issuer', _issuer, l10n.lmsIssuer),
              _field('client', _client, l10n.lmsClientId),
              _field('auth', _auth, l10n.lmsAuthUrl),
              _field('token', _token, l10n.lmsTokenUrl),
              _field('jwks', _jwks, l10n.lmsJwksUrl),
              TextField(
                key: const ValueKey('lms-deployments'),
                controller: _deployments,
                decoration: InputDecoration(
                  labelText: l10n.lmsDeployments,
                  helperText: l10n.lmsDeploymentsHint,
                ),
              ),
              SwitchListTile(
                contentPadding: EdgeInsets.zero,
                title: Text(l10n.lmsMatchEmail),
                subtitle: Text(l10n.lmsMatchEmailHint),
                value: _matchEmail,
                onChanged: (value) => setState(() => _matchEmail = value),
              ),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(l10n.actionCancel),
        ),
        FilledButton(
          key: const ValueKey('lms-save'),
          onPressed: _done,
          child: Text(l10n.actionSave),
        ),
      ],
    );
  }
}
