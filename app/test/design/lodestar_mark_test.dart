import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/design/components/lodestar_mark.dart';
import 'package:lodestar/design/components/lodestar_mark_painter.dart';
import 'package:material_ui/material_ui.dart';

void main() {
  testWidgets('mark keeps the requested square size', (tester) async {
    await tester.pumpWidget(
      const Center(
        child: LodestarMark(size: 40, color: Color(0xFF1E2650)),
      ),
    );
    expect(tester.getSize(find.byType(LodestarMark)), const Size(40, 40));
  });

  test('mark repaints only when the color changes', () {
    const ink = LodestarMarkPainter(color: Color(0xFF1E2650));
    expect(
      ink.shouldRepaint(const LodestarMarkPainter(color: Color(0xFF1E2650))),
      isFalse,
    );
    expect(
      ink.shouldRepaint(const LodestarMarkPainter(color: Color(0xFFF4EFE6))),
      isTrue,
    );
  });
}
