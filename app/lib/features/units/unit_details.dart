import 'unit.dart';

class UnitDetails {
  const UnitDetails({required this.unit, required this.path});

  factory UnitDetails.fromJson(Map<String, Object?> json) {
    return UnitDetails(
      unit: Unit.fromJson(json),
      path: [
        for (final item in json['path']! as List<Object?>)
          Unit.fromJson(item! as Map<String, Object?>),
      ],
    );
  }

  final Unit unit;
  final List<Unit> path;
}
