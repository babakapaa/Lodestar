import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../l10n/l10n_context.dart';
import 'course_material.dart';
import 'navigator_providers.dart';

class MaterialTile extends ConsumerWidget {
  const MaterialTile({
    required this.courseId,
    required this.material,
    super.key,
  });

  final String courseId;
  final CourseMaterial material;

  Future<void> _delete(BuildContext context, WidgetRef ref) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.materialDeleteConfirm(material.title)),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: Text(l10n.actionCancel),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: Text(l10n.materialDelete),
          ),
        ],
      ),
    );
    if (!(confirmed ?? false)) {
      return;
    }
    try {
      await ref.read(navigatorApiProvider).deleteMaterial(material.id);
      ref.invalidate(materialsProvider(courseId));
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Icon(
        material.ready ? Icons.task_alt : Icons.hourglass_top_outlined,
      ),
      title: Text(material.title),
      subtitle: Text(
        [
          if (material.reference.isNotEmpty) material.reference,
          material.ready
              ? l10n.materialReady(material.passages)
              : l10n.materialIndexing(material.indexed, material.passages),
        ].join('\n'),
      ),
      trailing: IconButton(
        tooltip: l10n.materialDelete,
        icon: const Icon(Icons.delete_outline),
        onPressed: () => _delete(context, ref),
      ),
    );
  }
}
