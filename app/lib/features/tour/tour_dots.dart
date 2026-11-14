import 'package:material_ui/material_ui.dart';

import '../../design/tokens/tokens_context.dart';

class TourDots extends StatelessWidget {
  const TourDots({required this.count, required this.current, super.key});

  final int count;
  final int current;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        for (var index = 0; index < count; index++)
          AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            margin: const EdgeInsets.symmetric(horizontal: 3),
            width: index == current ? 18 : 6,
            height: 6,
            decoration: BoxDecoration(
              color: index == current ? palette.accent : palette.border,
              borderRadius: BorderRadius.circular(3),
            ),
          ),
      ],
    );
  }
}
