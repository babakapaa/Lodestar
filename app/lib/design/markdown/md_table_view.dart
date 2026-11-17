import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';
import 'md_inline.dart';

class MdTableView extends StatelessWidget {
  const MdTableView({
    required this.rows,
    required this.base,
    required this.inline,
    super.key,
  });

  final List<List<String>> rows;
  final TextStyle base;
  final MdInline inline;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    final width = rows.fold(
      0,
      (max, row) => row.length > max ? row.length : max,
    );
    final head = base.copyWith(fontWeight: FontWeight.w700);
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Container(
        decoration: BoxDecoration(
          border: Border.all(color: palette.border),
          borderRadius: context.tokens.borderRadiusSmall,
        ),
        clipBehavior: Clip.antiAlias,
        child: Table(
          defaultColumnWidth: const IntrinsicColumnWidth(),
          border: TableBorder.symmetric(
            inside: BorderSide(color: palette.border),
          ),
          children: [
            for (var r = 0; r < rows.length; r++)
              TableRow(
                decoration: r == 0
                    ? BoxDecoration(color: palette.surface2)
                    : null,
                children: [
                  for (var c = 0; c < width; c++)
                    Padding(
                      padding: const EdgeInsets.symmetric(
                        horizontal: Spacing.m,
                        vertical: Spacing.s,
                      ),
                      child: Text.rich(
                        TextSpan(
                          children: inline.spans(
                            c < rows[r].length ? rows[r][c] : '',
                            r == 0 ? head : base,
                          ),
                        ),
                      ),
                    ),
                ],
              ),
          ],
        ),
      ),
    );
  }
}
