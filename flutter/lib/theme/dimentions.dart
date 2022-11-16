import 'package:flutter/material.dart';
import 'package:get/get.dart';

class KSize {
  static double getWidth(BuildContext context, width) {
    double _width =
        (((100 / 375) * width) / 100) * Get.width;
    return _width;
  }

  static double geHeight(BuildContext context, height) {
    double _height =
        (((100 / 812) * height) / 100) * Get.height;
    return _height;
  }
}