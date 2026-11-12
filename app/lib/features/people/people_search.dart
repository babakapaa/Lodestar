import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'people_providers.dart';
import 'person_hit.dart';

class PeopleSearch extends ConsumerStatefulWidget {
  const PeopleSearch({required this.onSelected, super.key});

  final ValueChanged<PersonHit> onSelected;

  @override
  ConsumerState<PeopleSearch> createState() => _PeopleSearchState();
}

class _PeopleSearchState extends ConsumerState<PeopleSearch> {
  static const _debounce = Duration(milliseconds: 300);

  Timer? _timer;
  String _query = '';

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  void _changed(String text) {
    _timer?.cancel();
    _timer = Timer(_debounce, () => setState(() => _query = text.trim()));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final results = _query.isEmpty
        ? const AsyncData<List<PersonHit>>([])
        : ref.watch(peopleSearchProvider(_query));
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TextField(
          autofocus: true,
          onChanged: _changed,
          decoration: InputDecoration(
            prefixIcon: const Icon(Icons.search),
            hintText: l10n.peopleSearchHint,
          ),
        ),
        const SizedBox(height: Spacing.m),
        switch (results) {
          AsyncData(:final value) when value.isEmpty && _query.isNotEmpty =>
            Text(l10n.peopleNotFound),
          AsyncData(:final value) => Column(
            children: [
              for (final person in value)
                ListTile(
                  contentPadding: EdgeInsets.zero,
                  leading: const Icon(Icons.person_outline),
                  title: Text(person.name),
                  subtitle: person.subtitle.isEmpty
                      ? null
                      : Text(person.subtitle),
                  onTap: () => widget.onSelected(person),
                ),
            ],
          ),
          AsyncError(:final error) => FormError(
            message: ApiErrorText.of(l10n, error),
          ),
          _ => const Center(child: CircularProgressIndicator()),
        },
      ],
    );
  }
}
