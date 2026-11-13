import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import 'chat_file.dart';

class ChatFileLink extends StatelessWidget {
  const ChatFileLink({required this.file, this.onTap, super.key});

  final ChatFile file;
  final VoidCallback? onTap;

  static String size(int bytes) {
    if (bytes < 1024) {
      return '$bytes B';
    }
    if (bytes < 1024 * 1024) {
      return '${(bytes / 1024).toStringAsFixed(1)} KB';
    }
    return '${(bytes / 1024 / 1024).toStringAsFixed(1)} MB';
  }

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    return InkWell(
      onTap: onTap,
      borderRadius: tokens.borderRadius,
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.insert_drive_file_outlined, color: tokens.palette.accent),
          const SizedBox(width: Spacing.xs),
          Flexible(child: Text(file.name, overflow: TextOverflow.ellipsis)),
          const SizedBox(width: Spacing.s),
          Text(
            size(file.size),
            style: Theme.of(context).textTheme.labelSmall
                ?.copyWith(color: tokens.palette.muted),
          ),
        ],
      ),
    );
  }
}
