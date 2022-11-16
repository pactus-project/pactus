import 'package:day_night_switcher/day_night_switcher.dart';
import 'package:flutter/material.dart';
import 'package:iconify_flutter/iconify_flutter.dart';
import 'package:iconify_flutter/icons/fa_solid.dart';
import 'package:iconify_flutter/icons/ic.dart';
import 'package:iconify_flutter/icons/ion.dart';
import 'package:iconify_flutter/icons/ph.dart';
import 'package:iconify_flutter/icons/zondicons.dart';
import 'package:pactus_app/theme/app_colors.dart';
import 'package:pactus_app/theme/dimentions.dart';
import 'package:pactus_app/widgets/default_background_color.dart';
import 'package:pactus_app/widgets/main_menu_button.dart';

class MainScreen extends StatelessWidget {
  const MainScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: DefaultBackgroundScreen(
          child: Container(
         height: MediaQuery.of(context).size.height,
            child: Column(
        children: [
            SizedBox(
              height: KSize.geHeight(context, 50),
            ),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: KSize.getWidth(context, 30)),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  DayNightSwitcherIcon(
                    
  isDarkModeEnabled: true,
  onStateChanged: (isDarkModeEnabled) {
    // setState(() {
    //   this.isDarkModeEnabled = isDarkModeEnabled;
    // });
  },
),
                  Container(
         width: KSize.getWidth(context, 60),
                    child: TextButton(
                    onPressed: (){},
                      child: Row(
                        
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Iconify(
                            FaSolid.language,
                            size: 25,
                            color: AppColors.white,
                          ),
                          SizedBox(width: KSize.getWidth(context, 3),),
                          Text('Englis(US)',style: Theme.of(context).textTheme.subtitle1,),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
      
            Padding(
              padding: EdgeInsets.symmetric(
                horizontal: MediaQuery.of(context).size.width / 4.5,
                // vertical: MediaQuery.of(context).size.height / 10
              ),
              child: Column(
                // mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Welcome to Pactus',
                    style: Theme.of(context).textTheme.headline5,
                  ),
                  SizedBox(
                    height: 30,
                  ),
                  MainMenuButton(
                      onPress: () {},
                      title: 'Create a new wallet',
                      icon: Ic.baseline_add_to_photos,
                      description:
                          'Choose this option if this is your first time using Pactus'),
                  Divider(
                    color: AppColors.white,
                    thickness: 0.25,
                    height: KSize.geHeight(context, 30),
                  ),
                  MainMenuButton(
                      onPress: () {},
                      title: 'Create a new wallet from hardware',
                      icon: Zondicons.hard_drive,
                      description:
                          'Connect your hardware wallet to create a new Pactus wallet'),
                  Divider(
                    color: AppColors.white,
                    thickness: 0.25,
                    height: KSize.geHeight(context, 30),
                  ),
                  MainMenuButton(
                      onPress: () {},
                      title: 'Open a wallet from file',
                      icon: Ion.folder_open_sharp,
                      description:
                          'Input an existing . keys wallet from your computer'),
                  Divider(
                    color: AppColors.white,
                    thickness: 0.25,
                    height: KSize.geHeight(context, 30),
                  ),
                  MainMenuButton(
                      onPress: () {},
                      title: 'Restore wallet from keys or mnemonic seed',
                      icon: Ph.password_fill,
                      description:
                          'Enter your private keys or 25-word mnemonic seed to restore your wallet'),      SizedBox(
              height: KSize.geHeight(context,60),
            ),
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton(onPressed: (){}, child: Text('Change wallet mode'), style: ElevatedButton.styleFrom(
                padding: EdgeInsets.symmetric(horizontal: 12 ,vertical: 6),
                backgroundColor: AppColors.main_color,
                shadowColor: AppColors.grey1,
                disabledBackgroundColor: AppColors.grey1,
                surfaceTintColor: AppColors.grey1,
                foregroundColor: AppColors.grey1,
                disabledForegroundColor: AppColors.grey1),
                  ),
            ),
                SizedBox(
              height: KSize.geHeight(context,60),
            ),
                ],
              ),
            ),
           
        ],
      ),
          )),
    );
  }
}
