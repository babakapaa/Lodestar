import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'tour_controller.dart';

final tourSeenProvider = NotifierProvider<TourController, bool>(
  TourController.new,
);
