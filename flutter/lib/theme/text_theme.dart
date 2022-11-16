import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:pactus_app/theme/app_colors.dart';

ThemeData getAppThemeData() {
  return ThemeData(
      textButtonTheme: TextButtonThemeData(
        style: ElevatedButton.styleFrom(
            backgroundColor: Colors.transparent,
            shadowColor: AppColors.white,
            disabledBackgroundColor: AppColors.white,
            surfaceTintColor: AppColors.white,
            foregroundColor: AppColors.white,
            disabledForegroundColor: AppColors.white),
      ),
      textTheme: TextTheme(
        headline4: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 34,
          fontWeight: FontWeight.w500,
        ),
        headline2: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 34,
          fontWeight: FontWeight.w500,
        ),
        headline3: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 30,
          fontWeight: FontWeight.bold,
        ),
        headline5: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 26,
          // fontWeight: FontWeight.bold,
          color: AppColors.white
        ),
        headline6: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 24,
          fontWeight: FontWeight.w500,
        ),
        subtitle1: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 16,
             color: AppColors.white,
          fontWeight: FontWeight.normal,
        ),
        overline: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 10,
          fontWeight: FontWeight.normal,
        ),
        caption: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 12,
          fontWeight: FontWeight.normal,
          color: AppColors.grey1
        ),
        button: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 14,
          fontWeight: FontWeight.w500,
        ),
        bodyText1: TextStyle(
          fontFamily: 'Poppins',
          fontSize: 14,
        ),
        bodyText2: TextStyle(
          fontFamily: 'Poppins',
          color: AppColors.white,
          fontSize: 16,
          fontWeight: FontWeight.normal,
        ),
      ),
      fontFamily: 'Poppins');
}
