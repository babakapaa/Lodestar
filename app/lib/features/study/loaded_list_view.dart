import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/offline/loaded.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';

class LoadedListView<T> extends StatelessWidget {
  const LoadedListView({
    required this.value,
    required this.onRefresh,
    required this.builder,
    this.bottomInset = 0,
    this.maxWidth = 720,
    super.key,
  });

  final AsyncValue<Loaded<T>> value;
  final Future<void> Function() onRefresh;
  final List<Widget> Function(T value) builder;
  final double bottomInset;
  final double maxWidth;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final loaded = value.value;
    final error = value.error;
    final wide = MediaQuery.sizeOf(context).width >= 600;
    final side = wide ? Spacing.xxl : Spacing.l;
    final padding = EdgeInsets.fromLTRB(
      side,
      wide ? Spacing.xl : Spacing.l,
      side,
      Spacing.xl + bottomInset,
    );
    return Align(
      alignment: Alignment.topCenter,
      child: ConstrainedBox(
        constraints: BoxConstraints(maxWidth: maxWidth),
        child: loaded != null
            ? RefreshIndicator(
                onRefresh: onRefresh,
                child: ListView(
                  padding: padding,
                  children: [
                    if (loaded.stale) ...[
                      NoticeBanner(
                        icon: Icons.cloud_off_outlined,
                        message: l10n.studyStale,
                      ),
                      const SizedBox(height: Spacing.l),
                    ],
                    ...builder(loaded.value),
                  ],
                ),
              )
            : error != null
            ? ListView(
                padding: padding,
                children: [
                  FormError(message: ApiErrorText.of(l10n, error)),
                  const SizedBox(height: Spacing.l),
                  Align(
                    alignment: Alignment.centerLeft,
                    child: LodestarButton(
                      label: l10n.actionRetry,
                      icon: Icons.refresh,
                      variant: LodestarButtonVariant.secondary,
                      onPressed: onRefresh,
                    ),
                  ),
                ],
              )
            : const Center(child: CircularProgressIndicator()),
      ),
    );
  }
}
