import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/study/map_gap.dart';
import 'package:lodestar/features/study/map_hint.dart';
import 'package:lodestar/features/study/map_rules.dart';
import 'package:lodestar/features/study/map_source.dart';
import 'package:lodestar/features/study/search_map.dart';

void main() {
  const article = MapSource(
    sourceId: 's1',
    reference: 'Статья',
    kind: 'article',
    why: 'Свежие данные о мембране',
  );
  const site = MapSource(
    sourceId: 's2',
    reference: 'Сайт',
    kind: 'web',
    why: 'Обзор для начала',
  );
  const rules = MapRules(minSources: 2, minScholarly: 60);

  test('an empty map misses sources and path', () {
    expect(const SearchMap().gaps(MapRules.none), [
      MapGap.sources,
      MapGap.path,
    ]);
  });

  test('the scholarly share and the count follow the rules', () {
    const map = SearchMap(used: [article, site], path: 'Путь');
    expect(map.scholarlyShare, 50);
    expect(map.gaps(rules), [MapGap.scholarly]);
    expect(map.gaps(const MapRules(minSources: 3, minScholarly: 50)), [
      MapGap.sources,
    ]);
  });

  test('a source with a DOI counts as scholarly', () {
    const book = MapSource(sourceId: 's3', reference: 'Книга', doi: '10.1/x');
    expect(book.scholarly, isTrue);
    expect(site.scholarly, isFalse);
  });

  test('every source needs its reason', () {
    const map = SearchMap(
      used: [article],
      rejected: [MapSource(sourceId: 's4', reference: 'Старое')],
      path: 'Путь',
    );
    expect(map.gaps(MapRules.none), [MapGap.why]);
  });

  test('a short reason is flagged as vague', () {
    expect(MapHint.vague('Полезно'), isTrue);
    expect(MapHint.vague(''), isFalse);
    expect(MapHint.vague('Дал схему строения мембраны'), isFalse);
  });
}
