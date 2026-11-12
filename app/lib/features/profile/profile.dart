import 'contact.dart';

class Profile {
  const Profile({
    required this.displayName,
    required this.about,
    required this.contacts,
  });

  factory Profile.fromJson(Map<String, Object?> json) {
    return Profile(
      displayName: json['displayName']! as String,
      about: json['about']! as String,
      contacts: [
        for (final item in json['contacts']! as List<Object?>)
          Contact.fromJson(item! as Map<String, Object?>),
      ],
    );
  }

  final String displayName;
  final String about;
  final List<Contact> contacts;

  Map<String, Object?> toJson() => {
    'displayName': displayName,
    'about': about,
    'contacts': [for (final contact in contacts) contact.toJson()],
  };
}
