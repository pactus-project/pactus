import 'package:flutter/material.dart';
import 'package:pactus_app/theme/app_colors.dart';
class DefaultBackgroundScreen extends StatelessWidget {
final Widget child;

  const DefaultBackgroundScreen({required this.child});
  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Container(
     
        child: Center(child: Container(
          decoration: BoxDecoration(gradient: LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [
    AppColors.background_light_color1,
    AppColors.background_light_color2,
    AppColors.background_light_color3,
    AppColors.background_light_color4,
    AppColors.background_light_color5,
    AppColors.background_light_color6,
    AppColors.background_light_color7,
    AppColors.background_light_color8,
    AppColors.background_dark_color1,
    AppColors.background_dark_color2,
    AppColors.background_dark_color3,
    AppColors.background_dark_color4,
    AppColors.background_dark_color5,
    AppColors.background_dark_color6,
    AppColors.background_dark_color7,
                ],
              )
        
    ),
    child: child,
        ),),
      ),
    );
  }
}