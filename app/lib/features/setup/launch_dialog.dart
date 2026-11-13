import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'pilot_course.dart';
import 'setup_providers.dart';

class LaunchDialog extends ConsumerStatefulWidget {
  const LaunchDialog({super.key});

  static Future<List<String>?> show(BuildContext context) =>
      showDialog<List<String>>(
        context: context,
        builder: (context) => const LaunchDialog(),
      );

  @override
  ConsumerState<LaunchDialog> createState() => _LaunchDialogState();
}

class _LaunchDialogState extends ConsumerState<LaunchDialog> {
  bool _pilot = false;
  final Set<String> _chosen = {};

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final courses = ref.watch(pilotCoursesProvider);
    return AlertDialog(
      title: Text(l10n.setupLaunchAction),
      content: SizedBox(
        width: 420,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              RadioGroup<bool>(
                groupValue: _pilot,
                onChanged: (v) => setState(() => _pilot = v ?? false),
                child: Column(
                  children: [
                    RadioListTile<bool>(
                      key: const ValueKey('launch-all'),
                      value: false,
                      title: Text(l10n.setupLaunchAll),
                    ),
                    RadioListTile<bool>(
                      key: const ValueKey('launch-pilot'),
                      value: true,
                      title: Text(l10n.setupLaunchPilot),
                      subtitle: Text(l10n.setupLaunchPilotHint),
                    ),
                  ],
                ),
              ),
              if (_pilot)
                ...switch (courses) {
                  AsyncData(:final value) when value.isEmpty => [
                    Text(l10n.setupPilotNoCourses),
                  ],
                  AsyncData(:final value) => [
                    for (final c in value) _courseTile(c),
                  ],
                  AsyncError() => [Text(l10n.errorGeneric)],
                  _ => [const LinearProgressIndicator()],
                },
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
          key: const ValueKey('launch-confirm'),
          onPressed: _pilot && _chosen.isEmpty
              ? null
              : () =>
                    Navigator.of(context)
                        .pop(_pilot ? [..._chosen] : <String>[]),
          child: Text(l10n.setupLaunchAction),
        ),
      ],
    );
  }

  Widget _courseTile(PilotCourse c) => CheckboxListTile(
    key: ValueKey('pilot-${c.id}'),
    value: _chosen.contains(c.id),
    title: Text(c.title),
    subtitle: c.code == null ? null : Text(c.code!),
    onChanged: (on) =>
        setState(() => on ?? false ? _chosen.add(c.id) : _chosen.remove(c.id)),
  );
}
