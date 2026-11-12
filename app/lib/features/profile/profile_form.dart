import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'contact_draft.dart';
import 'contact_kinds.dart';
import 'contact_row.dart';
import 'profile.dart';
import 'profile_providers.dart';

class ProfileForm extends ConsumerStatefulWidget {
  const ProfileForm({
    required this.profile,
    required this.fallbackName,
    super.key,
  });

  final Profile profile;
  final String fallbackName;

  @override
  ConsumerState<ProfileForm> createState() => _ProfileFormState();
}

class _ProfileFormState extends ConsumerState<ProfileForm> {
  late final _name = TextEditingController(text: widget.profile.displayName);
  late final _about = TextEditingController(text: widget.profile.about);
  late final List<ContactDraft> _contacts = [
    for (final contact in widget.profile.contacts) ContactDraft.from(contact),
  ];
  String? _error;
  bool _busy = false;

  @override
  void dispose() {
    _name.dispose();
    _about.dispose();
    for (final contact in _contacts) {
      contact.dispose();
    }
    super.dispose();
  }

  void _add() {
    setState(
      () => _contacts.add(
        ContactDraft(
          kind: ContactKinds.presets.first,
          value: '',
          public: false,
        ),
      ),
    );
  }

  void _remove(ContactDraft draft) {
    setState(() => _contacts.remove(draft));
    draft.dispose();
  }

  Future<void> _save() async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref
          .read(profileRepositoryProvider)
          .save(
            Profile(
              displayName: _name.text.trim(),
              about: _about.text.trim(),
              contacts: [
                for (final draft in _contacts)
                  if (draft.value.text.trim().isNotEmpty) draft.toContact(),
              ],
            ),
          );
      ref.invalidate(myProfileProvider);
      router.pop();
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
    return ListView(
      padding: const EdgeInsets.all(Spacing.l),
      children: [
        TextField(
          controller: _name,
          maxLength: 100,
          decoration: InputDecoration(
            labelText: l10n.profileName,
            hintText: widget.fallbackName,
            helperText: l10n.profileNameHint,
          ),
        ),
        const SizedBox(height: Spacing.m),
        TextField(
          controller: _about,
          maxLength: 2000,
          minLines: 3,
          maxLines: 8,
          decoration: InputDecoration(labelText: l10n.profileAbout),
        ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.contactsTitle,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              for (final draft in _contacts)
                ContactRow(
                  key: ObjectKey(draft),
                  draft: draft,
                  onChanged: () => setState(() {}),
                  onRemove: () => _remove(draft),
                ),
              Align(
                alignment: Alignment.centerLeft,
                child: LodestarButton(
                  label: l10n.contactAdd,
                  icon: Icons.add,
                  variant: LodestarButtonVariant.secondary,
                  onPressed: _contacts.length >= 50 ? null : _add,
                ),
              ),
            ],
          ),
        ),
        if (error != null) FormError(message: error),
        const SizedBox(height: Spacing.xl),
        LodestarButton(label: l10n.actionSave, onPressed: _busy ? null : _save),
      ],
    );
  }
}
