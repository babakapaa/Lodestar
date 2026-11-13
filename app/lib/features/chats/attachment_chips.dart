import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import 'chat_file.dart';

class AttachmentChips extends StatelessWidget {
  const AttachmentChips({
    required this.files,
    required this.onRemove,
    super.key,
  });

  final List<ChatFile> files;
  final ValueChanged<ChatFile> onRemove;

  @override
  Widget build(BuildContext context) => Wrap(
    spacing: Spacing.s,
    runSpacing: Spacing.xs,
    children: [
      for (final file in files)
        InputChip(
          avatar: const Icon(Icons.attach_file, size: 16),
          label: Text(file.name, overflow: TextOverflow.ellipsis),
          onDeleted: () => onRemove(file),
        ),
    ],
  );
}
