import 'apa_reference.dart';
import 'gost_reference.dart';
import 'source.dart';

enum ReferenceStyle {
  gost,
  apa;

  String list(List<Source> sources) => switch (this) {
    gost => GostReference.list(sources),
    apa => ApaReference.list(sources),
  };
}
