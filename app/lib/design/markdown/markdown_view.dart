import 'package:flutter_math_fork/flutter_math.dart';
import 'package:material_ui/material_ui.dart';

import '../tokens/spacing.dart';
import '../tokens/tokens_context.dart';
import '../typography/font_families.dart';
import 'md_block.dart';
import 'md_inline.dart';
import 'md_parser.dart';

class MarkdownView extends StatelessWidget {
  const MarkdownView(this.text, {this.style, super.key});

  static const _scales = [1.5, 1.3, 1.15, 1.05, 1.0, 1.0];

  final String text;
  final TextStyle? style;

  @override
  Widget build(BuildContext context) {
    final palette = context.palette;
    final base = (style ?? DefaultTextStyle.of(context).style).copyWith(
      height: 1.55,
    );
    final inline = MdInline(base, palette);
    final blocks = MdParser.parse(text);
    final children = <Widget>[];
    for (var i = 0; i < blocks.length; i++) {
      final block = blocks[i];
      final previous = i == 0 ? null : blocks[i - 1];
      if (i > 0) {
        final tight =
            block.kind == MdKind.item && previous?.kind == MdKind.item;
        children.add(SizedBox(height: tight ? Spacing.xs : Spacing.m));
      }
      children.add(_block(context, block, base, inline));
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      mainAxisSize: MainAxisSize.min,
      children: children,
    );
  }

  Widget _block(
    BuildContext context,
    MdBlock block,
    TextStyle base,
    MdInline inline,
  ) {
    final palette = context.palette;
    final size = base.fontSize ?? 15;
    switch (block.kind) {
      case MdKind.heading:
        final style = base.copyWith(
          fontSize: size * _scales[(block.level - 1).clamp(0, 5)],
          fontWeight: FontWeight.w700,
          height: 1.3,
        );
        return Padding(
          padding: const EdgeInsets.only(top: Spacing.xs),
          child: Text.rich(TextSpan(children: inline.spans(block.text, style))),
        );
      case MdKind.paragraph:
        return Text.rich(TextSpan(children: inline.spans(block.text)));
      case MdKind.item:
        final done = block.done;
        return Padding(
          padding: EdgeInsets.only(left: Spacing.l * block.level),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(
                width: 24,
                child: done != null
                    ? Padding(
                        padding: const EdgeInsets.only(top: 3),
                        child: Icon(
                          done
                              ? Icons.check_box
                              : Icons.check_box_outline_blank,
                          size: 18,
                          color: palette.accent,
                        ),
                      )
                    : Text(
                        block.marker,
                        style: base.copyWith(
                          color: palette.accent,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
              ),
              Expanded(
                child: Text.rich(
                  TextSpan(
                    children: inline.spans(
                      block.text,
                      done == true
                          ? base.copyWith(
                              color: palette.muted,
                              decoration: TextDecoration.lineThrough,
                            )
                          : null,
                    ),
                  ),
                ),
              ),
            ],
          ),
        );
      case MdKind.quote:
        return Container(
          padding: const EdgeInsets.fromLTRB(
            Spacing.m,
            Spacing.xs,
            0,
            Spacing.xs,
          ),
          decoration: BoxDecoration(
            border: Border(left: BorderSide(color: palette.accent, width: 3)),
          ),
          child: Text.rich(
            TextSpan(
              children: inline.spans(
                block.text,
                base.copyWith(
                  color: palette.muted,
                  fontStyle: FontStyle.italic,
                ),
              ),
            ),
          ),
        );
      case MdKind.rule:
        return Divider(color: palette.border, height: Spacing.l);
      case MdKind.code:
        return Container(
          padding: const EdgeInsets.all(Spacing.m),
          decoration: BoxDecoration(
            color: palette.surface2,
            borderRadius: context.tokens.borderRadiusSmall,
          ),
          child: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Text(
              block.text,
              style: base.copyWith(
                fontFamily: FontFamilies.jetBrainsMono,
                fontSize: size * 0.9,
                height: 1.45,
              ),
            ),
          ),
        );
      case MdKind.math:
        return Center(
          child: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Math.tex(
              block.text,
              mathStyle: MathStyle.display,
              textStyle: base.copyWith(fontSize: size * 1.1),
              onErrorFallback: (_) => Text(block.text, style: base),
            ),
          ),
        );
    }
  }
}
