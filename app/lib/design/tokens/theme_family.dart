import 'family_spec.dart';
import 'family_specs.dart';

enum ThemeFamily {
  paper(FamilySpecs.paper),
  graphite(FamilySpecs.graphite),
  lab(FamilySpecs.lab),
  campus(FamilySpecs.campus);

  const ThemeFamily(this.spec);

  final FamilySpec spec;
}
