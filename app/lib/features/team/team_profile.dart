class TeamProfile {
  const TeamProfile({
    required this.interests,
    required this.skills,
    required this.hours,
    required this.role,
  });

  factory TeamProfile.fromJson(Map<String, Object?> json) => TeamProfile(
    interests: json['interests'] as String? ?? '',
    skills: json['skills'] as String? ?? '',
    hours: json['hours'] as String? ?? '',
    role: json['role'] as String? ?? '',
  );

  final String interests;
  final String skills;
  final String hours;
  final String role;

  Map<String, Object?> toJson() => {
    'interests': interests,
    'skills': skills,
    'hours': hours,
    'role': role,
  };
}
