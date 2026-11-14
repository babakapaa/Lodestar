import 'package:material_ui/material_ui.dart';

import 'markdown_controller.dart';

class MarkdownText extends StatefulWidget {
  const MarkdownText(this.text, {this.style, super.key});

  final String text;
  final TextStyle? style;

  @override
  State<MarkdownText> createState() => _MarkdownTextState();
}

class _MarkdownTextState extends State<MarkdownText> {
  late final _controller = MarkdownController()..text = widget.text;

  @override
  void didUpdateWidget(MarkdownText oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.text != widget.text) {
      _controller.text = widget.text;
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SelectableText.rich(
      _controller.buildTextSpan(
        context: context,
        style: widget.style ?? DefaultTextStyle.of(context).style,
        withComposing: false,
      ),
    );
  }
}
