import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import 'tour_dialog.dart';
import 'tour_providers.dart';

class TourGate extends ConsumerStatefulWidget {
  const TourGate({required this.child, super.key});

  final Widget child;

  @override
  ConsumerState<TourGate> createState() => _TourGateState();
}

class _TourGateState extends ConsumerState<TourGate> {
  bool _showing = false;

  @override
  void initState() {
    super.initState();
    ref.listenManual(tourSeenProvider, (_, seen) {
      if (!seen) {
        WidgetsBinding.instance.addPostFrameCallback((_) => _show());
      }
    }, fireImmediately: true);
  }

  Future<void> _show() async {
    if (_showing || !mounted) {
      return;
    }
    _showing = true;
    await TourDialog.show(context);
    _showing = false;
    await ref.read(tourSeenProvider.notifier).finish();
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
