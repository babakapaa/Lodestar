import 'tester_person.dart';

class TesterInfo {
  const TesterInfo({required this.tenant, required this.people});

  factory TesterInfo.fromJson(Map<String, Object?> json) {
    return TesterInfo(
      tenant: json['tenant']! as String,
      people: [
        for (final item in json['people']! as List<Object?>)
          TesterPerson.fromJson(item! as Map<String, Object?>),
      ],
    );
  }

  final String tenant;
  final List<TesterPerson> people;
}
