import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'map_hint.dart';
import 'map_source.dart';

class MapSourceCard extends StatefulWidget {
  const MapSourceCard({
    required this.source,
    required this.rejected,
    required this.onChanged,
    required this.onMove,
    required this.onRemove,
    super.key,
  });

  final MapSource source;
  final bool rejected;
  final ValueChanged<MapSource> onChanged;
  final VoidCallback onMove;
  final VoidCallback onRemove;

  @override
  State<MapSourceCard> createState() => _MapSourceCardState();
}

class _MapSourceCardState extends State<MapSourceCard> {
  late final _found = TextEditingController(text: widget.source.found);
  late final _why = TextEditingController(text: widget.source.why);
  late final _taken = TextEditingController(text: widget.source.taken);

  @override
  void dispose() {
    _found.dispose();
    _why.dispose();
    _taken.dispose();
    super.dispose();
  }

  void _changed() {
    widget.onChanged(
      widget.source.copyWith(
        found: _found.text,
        why: _why.text,
        taken: _taken.text,
      ),
    );
  }

  Widget _field(
    TextEditingController controller,
    String label, {
    String? helper,
  }) => Padding(
    padding: const EdgeInsets.only(top: Spacing.s),
    child: TextField(
      controller: controller,
      minLines: 1,
      maxLines: 4,
      maxLength: 2000,
      buildCounter: (
        _, {
        required currentLength,
        required isFocused,
        maxLength,
      }) => null,
      decoration: InputDecoration(labelText: label, helperText: helper),
      onChanged: (_) => _changed(),
    ),
  );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = context.palette;
    final source = widget.source;
    final muted = Theme.of(context).textTheme.bodySmall!
        .copyWith(color: palette.muted);
    return LodestarCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(child: SelectableText(source.reference)),
              PopupMenuButton<bool>(
                onSelected: (move) =>
                    move ? widget.onMove() : widget.onRemove(),
                itemBuilder: (context) => [
                  PopupMenuItem(
                    value: true,
                    child: Text(
                      widget.rejected ? l10n.mapMoveUsed : l10n.mapMoveRejected,
                    ),
                  ),
                  PopupMenuItem(value: false, child: Text(l10n.mapRemove)),
                ],
              ),
            ],
          ),
          if (source.scholarly) Text(l10n.mapScholarly, style: muted),
          _field(_found, l10n.mapFound),
          _field(
            _why,
            widget.rejected ? l10n.mapWhyRejected : l10n.mapWhyUsed,
            helper: MapHint.vague(_why.text) ? l10n.mapVague : null,
          ),
          if (!widget.rejected) _field(_taken, l10n.mapTaken),
        ],
      ),
    );
  }
}
