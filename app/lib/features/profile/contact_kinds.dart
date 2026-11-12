import '../../l10n/generated/app_localizations.dart';

abstract final class ContactKinds {
  static const presets = [
    'telegram',
    'whatsapp',
    'phone',
    'email',
    'website',
    'other',
  ];

  static String label(AppLocalizations l10n, String kind) => switch (kind) {
    'telegram' => 'Telegram',
    'whatsapp' => 'WhatsApp',
    'phone' => l10n.contactPhone,
    'email' => l10n.contactEmail,
    'website' => l10n.contactWebsite,
    'other' => l10n.contactOther,
    _ => kind,
  };
}
