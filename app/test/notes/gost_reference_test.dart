import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/notes/gost_reference.dart';
import 'package:lodestar/features/notes/source.dart';
import 'package:lodestar/features/notes/source_kind.dart';

void main() {
  test('a book with one author', () {
    const book = Source(
      id: '1',
      kind: SourceKind.book,
      title: 'Основы цитологии',
      authors: ['Иванов, Алексей Борисович'],
      city: 'Алматы',
      publisher: 'Қазақ университеті',
      year: '2020',
      isbn: '978-601-04-0000-0',
    );
    expect(
      GostReference.format(book),
      'Иванов, А. Б. Основы цитологии / А. Б. Иванов. – '
      'Алматы : Қазақ университеті, 2020. – ISBN 978-601-04-0000-0.',
    );
  });

  test('a journal article with a DOI', () {
    const article = Source(
      id: '2',
      kind: SourceKind.article,
      title: 'Строение клетки',
      authors: ['Иванов, А. Б.', 'Ким, Е.'],
      container: 'Вестник биологии',
      year: '2021',
      volume: '12',
      issue: '3',
      pages: '45–67',
      doi: '10.1000/xyz',
    );
    expect(
      GostReference.format(article),
      'Иванов, А. Б. Строение клетки / А. Б. Иванов, Е. Ким // Вестник биологии. – '
      '2021. – Т. 12, № 3. – С. 45–67. – DOI 10.1000/xyz.',
    );
  });

  test('an English article with five authors', () {
    const article = Source(
      id: '3',
      kind: SourceKind.article,
      title: 'Cell division',
      authors: ['Smith, John', 'Lee, Ann', 'Kim, Bo', 'Park, Ji', 'Doe, Jane'],
      container: 'Nature',
      year: '2019',
      volume: '5',
      pages: '1–9',
    );
    expect(
      GostReference.format(article),
      'Cell division / J. Smith, A. Lee, B. Kim [et al.] // Nature. – 2019. – '
      'Vol. 5. – P. 1–9.',
    );
  });

  test('a web page with the access date', () {
    const page = Source(
      id: '4',
      kind: SourceKind.web,
      title: 'Митоз',
      container: 'Энциклопедия',
      url: 'https://example.org/mitoz',
      accessedAt: '2026-11-12',
    );
    expect(
      GostReference.format(page),
      'Митоз : электронный ресурс // Энциклопедия. – '
      'URL: https://example.org/mitoz (дата обращения: 12.11.2026).',
    );
  });

  test('the list is numbered and sorted', () {
    expect(
      GostReference.list(const [
        Source(id: 'a', kind: SourceKind.book, title: 'Ятаган'),
        Source(id: 'b', kind: SourceKind.book, title: 'Азбука'),
      ]),
      '1. Азбука.\n2. Ятаган.',
    );
  });
}
