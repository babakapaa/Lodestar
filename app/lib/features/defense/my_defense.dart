import '../study/defense_format.dart';
import 'defense.dart';
import 'defense_slot.dart';
import 'defense_status.dart';

class MyDefense {
  const MyDefense({
    required this.format,
    required this.required,
    required this.recorded,
    required this.status,
    required this.slot,
    required this.protocol,
  });

  factory MyDefense.fromJson(Map<String, Object?> json) {
    final slot = json['slot'];
    final protocol = json['protocol'];
    return MyDefense(
      format: DefenseFormat.parse(json['format']),
      required: json['required'] as bool? ?? false,
      recorded: json['recorded'] as bool? ?? false,
      status: DefenseStatus.parse(json['status']),
      slot: slot is Map<String, Object?> ? DefenseSlot.fromJson(slot) : null,
      protocol: protocol is Map<String, Object?>
          ? Defense.fromJson(protocol)
          : null,
    );
  }

  final DefenseFormat format;
  final bool required;
  final bool recorded;
  final DefenseStatus status;
  final DefenseSlot? slot;
  final Defense? protocol;
}
