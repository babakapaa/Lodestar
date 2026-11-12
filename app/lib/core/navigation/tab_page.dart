import 'package:material_ui/material_ui.dart';

import 'layout_size.dart';
import 'notifications_button.dart';
import 'profile_button.dart';

class TabPage extends StatelessWidget {
  const TabPage({
    required this.title,
    required this.child,
    this.floatingActionButton,
    super.key,
  });

  final String title;
  final Widget child;
  final Widget? floatingActionButton;

  @override
  Widget build(BuildContext context) {
    final compact =
        LayoutSize.of(MediaQuery.sizeOf(context).width) == LayoutSize.compact;
    return Scaffold(
      appBar: AppBar(
        title: Text(title),
        actions: [
          const NotificationsButton(),
          if (compact) const ProfileButton(showLabel: false),
        ],
      ),
      floatingActionButton: floatingActionButton,
      body: SafeArea(top: false, child: child),
    );
  }
}
