import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../defense/my_defense_section.dart';
import '../sandbox/sandbox_work_tile.dart';
import '../team/team_work_section.dart';
import 'assignment_details.dart';
import 'attached_file_tile.dart';
import 'class_work_section.dart';
import 'course_role.dart';
import 'defense_format.dart';
import 'loaded_list_view.dart';
import 'my_work_section.dart';
import 'navigator_tile.dart';
import 'search_map_tile.dart';
import 'study_labels_context.dart';
import 'study_providers.dart';
import 'submit_form.dart';

class AssignmentScreen extends ConsumerWidget {
  const AssignmentScreen({required this.assignmentId, super.key});

  final String assignmentId;

  Future<void> _edit(BuildContext context) async {
    final deleted = await context.push<bool>(
      AppRoutes.assignmentEdit(assignmentId),
    );
    if ((deleted ?? false) && context.mounted) {
      context.pop();
    }
  }

  Future<void> _refresh(WidgetRef ref) {
    ref.invalidate(assignmentWorkProvider(assignmentId));
    return ref.refresh(assignmentProvider(assignmentId).future);
  }

  List<Widget> _body(BuildContext context, AssignmentDetails details) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final assignment = details.assignment;
    final textTheme = Theme.of(context).textTheme;
    return [
      Text(assignment.title, style: textTheme.headlineSmall),
      const SizedBox(height: Spacing.s),
      Text(
        [
          labels.due(assignment.dueAt),
          labels.grading(assignment.grading),
          if (!assignment.published) l10n.assignmentDraft,
        ].join(', '),
        style: textTheme.bodyMedium!.copyWith(color: context.palette.muted),
      ),
      if (assignment.description.isNotEmpty) ...[
        const SizedBox(height: Spacing.l),
        SelectableText(assignment.description),
      ],
      if (details.files.isNotEmpty) ...[
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.assignmentFiles,
          child: Column(
            children: [
              for (final file in details.files) AttachedFileTile(file: file),
            ],
          ),
        ),
      ],
      const SizedBox(height: Spacing.xl),
      if (details.role == CourseRole.student) ...[
        MyWorkSection(details: details),
        const SizedBox(height: Spacing.xl),
        NavigatorTile(
          icon: Icons.explore_outlined,
          title: l10n.navigatorAsk,
          route: AppRoutes.newDialog(
            scopeKind: 'assignment',
            scopeId: assignment.id,
          ),
        ),
        const SizedBox(height: Spacing.l),
        SearchMapTile(assignment: assignment),
        const SizedBox(height: Spacing.l),
        SandboxWorkTile(
          kind: 'assignment',
          id: assignment.id,
          title: assignment.title,
        ),
        if (assignment.team == null)
          SubmitForm(details: details)
        else
          TeamWorkSection(details: details),
        if (assignment.defense != DefenseFormat.none) ...[
          const SizedBox(height: Spacing.xl),
          MyDefenseSection(assignmentId: assignment.id),
        ],
      ] else ...[
        if (assignment.team != null) ...[
          NavigatorTile(
            key: const ValueKey('teams'),
            icon: Icons.groups_outlined,
            title: l10n.teamsTitle,
            route: AppRoutes.teamScope(
              TeamWorkSection.scopeKind,
              assignment.id,
            ),
          ),
          const SizedBox(height: Spacing.l),
        ],
        if (assignment.defense != DefenseFormat.none) ...[
          NavigatorTile(
            key: const ValueKey('defenses'),
            icon: Icons.record_voice_over_outlined,
            title: l10n.defensesTitle,
            route: AppRoutes.defenses(assignment.id),
          ),
          const SizedBox(height: Spacing.l),
        ],
        ClassWorkSection(assignment: assignment),
      ],
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final details = ref.watch(assignmentProvider(assignmentId));
    final loaded = details.value?.value;
    return Scaffold(
      appBar: AppBar(
        title: Text(loaded?.courseTitle ?? ''),
        actions: [
          if (loaded != null && loaded.role.teaches)
            IconButton(
              tooltip: l10n.assignmentEdit,
              icon: const Icon(Icons.edit_outlined),
              onPressed: () => _edit(context),
            ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: LoadedListView(
          value: details,
          onRefresh: () => _refresh(ref),
          builder: (details) => _body(context, details),
        ),
      ),
    );
  }
}
