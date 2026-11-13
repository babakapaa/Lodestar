import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'team_profile.dart';
import 'team_providers.dart';

class TeamProfileForm extends ConsumerStatefulWidget {
  const TeamProfileForm({required this.profile, super.key});

  final TeamProfile profile;

  @override
  ConsumerState<TeamProfileForm> createState() => _TeamProfileFormState();
}

class _TeamProfileFormState extends ConsumerState<TeamProfileForm> {
  late final _interests = TextEditingController(text: widget.profile.interests);
  late final _skills = TextEditingController(text: widget.profile.skills);
  late final _hours = TextEditingController(text: widget.profile.hours);
  late final _role = TextEditingController(text: widget.profile.role);
  String? _error;
  bool _busy = false;
  bool _saved = false;

  @override
  void dispose() {
    _interests.dispose();
    _skills.dispose();
    _hours.dispose();
    _role.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    final l10n = context.l10n;
    setState(() {
      _busy = true;
      _error = null;
      _saved = false;
    });
    try {
      await ref
          .read(teamApiProvider)
          .saveProfile(
            TeamProfile(
              interests: _interests.text.trim(),
              skills: _skills.text.trim(),
              hours: _hours.text.trim(),
              role: _role.text.trim(),
            ),
          );
      setState(() => _saved = true);
    } on ApiException catch (error) {
      setState(() => _error = ApiErrorText.of(l10n, error));
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  Widget _field(String key, TextEditingController controller, String label) =>
      Padding(
        padding: const EdgeInsets.only(bottom: Spacing.l),
        child: TextField(
          key: ValueKey(key),
          controller: controller,
          minLines: 1,
          maxLines: 4,
          decoration: InputDecoration(labelText: label),
        ),
      );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(
          l10n.teamProfileAbout,
          style: Theme.of(context).textTheme.bodyMedium!
              .copyWith(color: context.palette.muted),
        ),
        const SizedBox(height: Spacing.l),
        _field('profile-interests', _interests, l10n.teamProfileInterests),
        _field('profile-skills', _skills, l10n.teamProfileSkills),
        _field('profile-hours', _hours, l10n.teamProfileHours),
        _field('profile-role', _role, l10n.teamProfileRole),
        if (error != null) FormError(message: error),
        if (_saved) Text(l10n.teamProfileSaved),
        const SizedBox(height: Spacing.l),
        LodestarButton(label: l10n.actionSave, onPressed: _busy ? null : _save),
      ],
    );
  }
}
