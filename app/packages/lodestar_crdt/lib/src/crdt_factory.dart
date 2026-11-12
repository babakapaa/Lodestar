import 'crdt_text.dart';
import 'native_crdt_text.dart';

CrdtText createCrdtText(int clientId) => NativeCrdtText(clientId);

Future<void> prepareCrdt() async {}
