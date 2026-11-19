import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/core/sources/apa_reference.dart';
import 'package:lodestar/core/sources/source.dart';
import 'package:lodestar/core/sources/source_kind.dart';

void main() {
  test('a journal article with a DOI', () {
    const article = Source(
      id: '1',
      kind: SourceKind.article,
      title: 'Photosynthesis: two phases',
      authors: ['Smith, Anna Maria', 'Jones, Bob'],
      container: 'Plant Biology',
      year: '2021',
      volume: '12',
      issue: '3',
      pages: '45–67',
      doi: '10.1000/xyz',
    );
    expect(
      ApaReference.format(article),
      'Smith, A. M., & Jones, B. (2021). Photosynthesis: two phases. '
      'Plant Biology, 12(3), 45–67. https://doi.org/10.1000/xyz',
    );
  });

  test('a book without authors or year', () {
    const book = Source(
      id: '2',
      kind: SourceKind.book,
      title: 'Cell biology.',
      publisher: 'Campus Press',
    );
    expect(ApaReference.format(book), 'Cell biology. (n.d.). Campus Press.');
  });

  test('a chapter and a thesis', () {
    const chapter = Source(
      id: '3',
      kind: SourceKind.chapter,
      title: 'Light reactions',
      authors: ['Kim, E.'],
      container: 'Plant physiology',
      publisher: 'Springer',
      year: '2019',
      pages: '10–20',
    );
    const thesis = Source(
      id: '4',
      kind: SourceKind.thesis,
      title: 'Calvin cycle in algae',
      authors: ['Lee, Ann'],
      publisher: 'KIMEP University',
      year: '2024',
      url: 'https://example.org/t',
    );
    expect(
      ApaReference.format(chapter),
      'Kim, E. (2019). Light reactions. In Plant physiology (pp. 10–20). '
      'Springer.',
    );
    expect(
      ApaReference.format(thesis),
      'Lee, A. (2024). Calvin cycle in algae [Thesis, KIMEP University]. '
      'https://example.org/t',
    );
  });

  test('the list is sorted by family name without numbers', () {
    const a = Source(
      id: 'a',
      kind: SourceKind.other,
      title: 'B',
      authors: ['Zhou, Li'],
    );
    const b = Source(
      id: 'b',
      kind: SourceKind.other,
      title: 'A',
      authors: ['Abay, Kunanbaiuly'],
    );
    expect(
      ApaReference.list([a, b]),
      'Abay, K. (n.d.). A.\nZhou, L. (n.d.). B.',
    );
  });
}
