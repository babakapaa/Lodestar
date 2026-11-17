import 'family_spec.dart';
import 'family_specs.dart';

enum ThemeFamily {
  campus(FamilySpecs.campus),
  paper(FamilySpecs.paper),
  graphite(FamilySpecs.graphite),
  lab(FamilySpecs.lab);

  const ThemeFamily(this.spec);

  final FamilySpec spec;
}
