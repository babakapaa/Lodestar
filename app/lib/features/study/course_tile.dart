import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../l10n/l10n_context.dart';
import 'course.dart';
import 'study_labels_context.dart';
import '../../design/components/tile_group.dart';

class CourseTile extends StatelessWidget {
  const CourseTile({required this.course, super.key});

  final Course course;

  @override
  Widget build(BuildContext context) {
    final code = course.code;
    return ListTile(
      contentPadding: TileGroup.tilePadding,
      leading: const Icon(Icons.school_outlined),
      title: Text(course.title),
      subtitle: Text(
        [
          if (code != null && code.isNotEmpty) code,
          context.studyLabels.role(course.role),
          if (course.archived) context.l10n.studyArchived,
        ].join(', '),
      ),
      trailing: const Icon(Icons.chevron_right),
      onTap: () => context.push(AppRoutes.course(course.id)),
    );
  }
}
