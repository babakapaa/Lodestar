import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../l10n/l10n_context.dart';
import '../notes/notes_providers.dart';
import 'article.dart';

class ArticleTile extends ConsumerStatefulWidget {
  const ArticleTile({required this.article, super.key});

  final Article article;

  @override
  ConsumerState<ArticleTile> createState() => _ArticleTileState();
}

class _ArticleTileState extends ConsumerState<ArticleTile> {
  bool _busy = false;
  bool _saved = false;

  Future<void> _save() async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    setState(() => _busy = true);
    try {
      final api = ref.read(sourcesApiProvider);
      await api.save(await api.lookup(doi: widget.article.doi));
      ref.invalidate(sourcesProvider);
      setState(() => _saved = true);
      messenger.showSnackBar(
        SnackBar(content: Text(l10n.navigatorArticleSaved)),
      );
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final article = widget.article;
    final byline = article.byline;
    return ListTile(
      dense: true,
      contentPadding: EdgeInsets.zero,
      leading: const Icon(Icons.article_outlined),
      title: Text(article.title),
      subtitle: Text(
        [if (byline.isNotEmpty) byline, 'DOI: ${article.doi}'].join('\n'),
        maxLines: 3,
        overflow: TextOverflow.ellipsis,
      ),
      trailing: IconButton(
        tooltip: context.l10n.navigatorArticleSave,
        onPressed: _busy || _saved ? null : _save,
        icon: Icon(_saved ? Icons.check : Icons.bookmark_add_outlined),
      ),
    );
  }
}
