import 'package:flutter/material.dart';
import 'package:iconify_flutter/iconify_flutter.dart';
import 'package:iconify_flutter/icons/ic.dart';
import 'package:pactus_app/theme/app_colors.dart';
import 'package:pactus_app/theme/dimentions.dart';

class MainMenuButton extends StatelessWidget {
 final VoidCallback onPress;
 final String title;
 final String icon;
 final String description;

  const MainMenuButton({super.key, required this.onPress, required this.title, required this.icon, required this.description});


  @override
  Widget build(BuildContext context) {
    return Column(
         crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        TextButton(
          
            style: Theme.of(context).textButtonTheme.style,
            onPressed: onPress,
            child: Padding(
              padding: EdgeInsets.symmetric(
                vertical: KSize.geHeight(context, 5)
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Iconify(
                   icon,
                    size: 35,
                    color: AppColors.main_color,
                  ),
                  SizedBox(
                    width: KSize.getWidth(context, 15),
                  ),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [  

                      Text(title,
                          style: Theme.of(context).textTheme.bodyText2),
                       SizedBox(
                    height: KSize.geHeight(context, 5),
                  ), 
                      Text(
                         description,
                          style: Theme.of(context).textTheme.caption),
                          
                    ],
                  ),
                ],
              ),
            )),
       
      ],
    );
  }
}