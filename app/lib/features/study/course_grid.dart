import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import 'course.dart';
import 'course_card.dart';

class CourseGrid extends StatelessWidget {
  const CourseGrid({required this.courses, super.key});

  final List<Course> courses;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final width = constraints.maxWidth;
        final columns = width >= 820 ? 3 : (width >= 520 ? 2 : 1);
        final cell = (width - Spacing.l * (columns - 1)) / columns;
        return Wrap(
          spacing: Spacing.l,
          runSpacing: Spacing.l,
          children: [
            for (final course in courses)
              SizedBox(
                width: cell,
                child: CourseCard(key: ValueKey(course.id), course: course),
              ),
          ],
        );
      },
    );
  }
}
