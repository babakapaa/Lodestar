import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_tile.dart';
import 'course_preview_providers.dart';
import 'study_labels_context.dart';

class CoursePreviewScreen extends ConsumerStatefulWidget {
  const CoursePreviewScreen({required this.courseId, super.key});

  final String courseId;

  @override
  ConsumerState<CoursePreviewScreen> createState() =>
      _CoursePreviewScreenState();
}

class _CoursePreviewScreenState extends ConsumerState<CoursePreviewScreen> {
  String? _studentId;
  DateTime _day = DateUtils.dateOnly(DateTime.now());

  DateTime get _at => _day.add(const Duration(hours: 23, minutes: 59));

  Future<void> _pickDay() async {
    final day = await showDatePicker(
      context: context,
      initialDate: _day,
      firstDate: _day.subtract(const Duration(days: 365)),
      lastDate: _day.add(const Duration(days: 365)),
    );
    if (day != null) {
      setState(() => _day = day);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final students =
        ref.watch(courseStudentsProvider(widget.courseId)).value ?? const [];
    final query = (courseId: widget.courseId, studentId: _studentId, at: _at);
    return Scaffold(
      appBar: AppBar(title: Text(l10n.previewTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(coursePreviewProvider(query)),
          onRefresh: () => ref.refresh(coursePreviewProvider(query).future),
          builder: (items) => [
            DropdownButtonFormField<String?>(
              key: const ValueKey('preview-student'),
              initialValue: _studentId,
              isExpanded: true,
              decoration: InputDecoration(labelText: l10n.previewStudent),
              items: [
                DropdownMenuItem(value: null, child: Text(l10n.previewAnyone)),
                for (final s in students)
                  DropdownMenuItem(value: s.userId, child: Text(s.label)),
              ],
              onChanged: (id) => setState(() => _studentId = id),
            ),
            const SizedBox(height: Spacing.s),
            ListTile(
              key: const ValueKey('preview-date'),
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.event_outlined),
              title: Text(l10n.previewDate),
              subtitle: Text(labels.date(_day)),
              onTap: _pickDay,
            ),
            const Divider(),
            if (items.isEmpty)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: Spacing.l),
                child: Text(l10n.previewEmpty),
              ),
            for (final item in items)
              AssignmentTile(
                key: ValueKey('preview-${item.assignment.id}'),
                item: item,
                teaches: false,
                at: _at,
              ),
          ],
        ),
      ),
    );
  }
}
