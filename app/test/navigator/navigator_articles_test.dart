import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/navigator/navigator_message.dart';

void main() {
  test('a reply carries OpenAlex articles', () {
    final message = NavigatorMessage.fromJson({
      'seq': 2,
      'role': 'navigator',
      'text': 'В материалах курса нет ответа.',
      'citations': <Object?>[],
      'articles': [
        {
          'title': 'Membrane transport',
          'authors': ['A. Smith', 'B. Lee'],
          'year': 2021,
          'venue': 'Cell',
          'doi': '10.1/abc',
          'url': 'https://doi.org/10.1/abc',
        },
      ],
    });
    final article = message.articles.single;
    expect(article.doi, '10.1/abc');
    expect(article.byline, 'A. Smith, B. Lee. 2021. Cell');
  });
}
