import 'package:material_ui/material_ui.dart';

import '../../l10n/l10n_context.dart';
import 'article.dart';
import 'article_tile.dart';

class ArticleList extends StatelessWidget {
  const ArticleList({required this.articles, super.key});

  final List<Article> articles;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          context.l10n.navigatorArticles,
          style: Theme.of(context).textTheme.labelLarge,
        ),
        for (final article in articles) ArticleTile(article: article),
      ],
    );
  }
}
