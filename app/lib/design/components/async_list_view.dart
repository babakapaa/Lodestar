import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../l10n/l10n_context.dart';
import '../tokens/spacing.dart';
import 'form_error.dart';

class AsyncListView<T> extends StatelessWidget {
  const AsyncListView({
    required this.value,
    required this.onRefresh,
    required this.builder,
    super.key,
  });

  final AsyncValue<T> value;
  final Future<void> Function() onRefresh;
  final List<Widget> Function(T value) builder;

  @override
  Widget build(BuildContext context) {
    final loaded = value.value;
    final error = value.error;
    return Align(
      alignment: Alignment.topCenter,
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 720),
        child: RefreshIndicator(
          onRefresh: onRefresh,
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              if (loaded != null)
                ...builder(loaded)
              else if (error != null)
                FormError(message: ApiErrorText.of(context.l10n, error))
              else
                const Center(child: CircularProgressIndicator()),
            ],
          ),
        ),
      ),
    );
  }
}
