import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';

class ShareBar extends StatelessWidget {
  const ShareBar({
    required this.label,
    required this.value,
    required this.share,
    super.key,
  });

  final String label;
  final String value;
  final double share;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: Spacing.xs),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(child: Text(label)),
              Text(value),
            ],
          ),
          const SizedBox(height: Spacing.xs),
          LinearProgressIndicator(value: share.clamp(0, 1)),
        ],
      ),
    );
  }
}
