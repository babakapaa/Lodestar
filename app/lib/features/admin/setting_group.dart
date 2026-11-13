enum SettingGroup {
  general,
  privacy,
  study,
  navigator,
  attendance,
  sandbox,
  other;

  static const _privacy = {
    'core.audit_retention_days',
    'core.self_deletion',
    'core.deletion_grace_days',
    'chat.retention_days',
    'study.defense_recording',
  };

  static SettingGroup of(String key) {
    if (_privacy.contains(key)) {
      return privacy;
    }
    return switch (key.split('.').first) {
      'core' => general,
      'study' || 'league' || 'team' => study,
      'navigator' => navigator,
      'attendance' => attendance,
      'sandbox' => sandbox,
      _ => other,
    };
  }
}
