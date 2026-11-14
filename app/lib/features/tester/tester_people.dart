import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'tester_person.dart';
import 'tester_providers.dart';

class TesterPeople extends ConsumerStatefulWidget {
  const TesterPeople({required this.people, this.current, super.key});

  final List<TesterPerson> people;
  final String? current;

  @override
  ConsumerState<TesterPeople> createState() => _TesterPeopleState();
}

class _TesterPeopleState extends ConsumerState<TesterPeople> {
  bool _busy = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) async {
      final next = ref.read(testerNextProvider.notifier).take();
      if (next == null || !mounted) {
        return;
      }
      await _shown();
      if (mounted) {
        await _signInAs(next);
      }
    });
  }

  Future<void> _shown() {
    final animation = ModalRoute.of(context)?.animation;
    if (animation == null || animation.status == AnimationStatus.completed) {
      return Future.value();
    }
    final shown = Completer<void>();
    void listener(AnimationStatus status) {
      if (status == AnimationStatus.completed) {
        animation.removeStatusListener(listener);
        shown.complete();
      }
    }

    animation.addStatusListener(listener);
    return shown.future;
  }

  Future<void> _signInAs(TesterPerson person) async {
    final l10n = context.l10n;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref.read(testerSwitchProvider).signInAs(person);
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
    final error = _error;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        for (final person in widget.people)
          Padding(
            padding: const EdgeInsets.only(bottom: Spacing.s),
            child: LodestarButton(
              label: context.l10n.testerSignInAs(person.name, person.title),
              icon: person.email == widget.current
                  ? Icons.check
                  : Icons.person_outline,
              variant: LodestarButtonVariant.secondary,
              onPressed: _busy || person.email == widget.current
                  ? null
                  : () => _signInAs(person),
            ),
          ),
        if (error != null) FormError(message: error),
      ],
    );
  }
}
