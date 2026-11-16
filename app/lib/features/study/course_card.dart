import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'course.dart';
import 'course_cover.dart';
import 'study_labels_context.dart';

class CourseCard extends StatelessWidget {
  const CourseCard({required this.course, super.key});

  final Course course;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    final textTheme = Theme.of(context).textTheme;
    final code = course.code;
    final colors = CourseCover.of(course.id + course.title, palette.isDark);
    return DecoratedBox(
      decoration: BoxDecoration(
        color: palette.surface,
        borderRadius: tokens.borderRadius,
        border: Border.all(color: palette.border),
        boxShadow: palette.raised,
      ),
      child: ClipRRect(
        borderRadius: tokens.borderRadius,
        child: Material(
          type: MaterialType.transparency,
          child: InkWell(
            onTap: () => context.push(AppRoutes.course(course.id)),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Container(
                  height: 88,
                  padding: const EdgeInsets.all(Spacing.l),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                      colors: colors,
                    ),
                  ),
                  alignment: Alignment.topLeft,
                  child: Row(
                    children: [
                      if (code != null && code.isNotEmpty)
                        DecoratedBox(
                          decoration: BoxDecoration(
                            color: const Color(0x33000000),
                            borderRadius: BorderRadius.circular(6),
                          ),
                          child: Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: Spacing.s,
                              vertical: 2,
                            ),
                            child: Text(
                              code,
                              style: textTheme.labelSmall!.copyWith(
                                color: Colors.white,
                              ),
                            ),
                          ),
                        ),
                      const Spacer(),
                      const Icon(
                        Icons.school_outlined,
                        color: Colors.white,
                        size: 22,
                      ),
                    ],
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(Spacing.l),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        course.title,
                        style: textTheme.titleMedium,
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: Spacing.xs),
                      Text(
                        [
                          context.studyLabels.role(course.role),
                          if (course.archived) context.l10n.studyArchived,
                        ].join(', '),
                        style: textTheme.bodySmall,
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
