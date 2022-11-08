<!--suppress ALL -->
<img src="assets/arkwaifu_phantom@0.25x.png" alt="logo" align="right" height="224" width="224"/>

# ArkWaifu (arkwaifu)

[![](https://pkg.go.dev/badge/github.com/flandiayingman/arkwaifu.svg)](https://pkg.go.dev/github.com/flandiayingman/arkwaifu)
![](https://img.shields.io/github/license/FlandiaYingman/arkwaifu?style=flat-square)
![](https://img.shields.io/github/last-commit/FlandiaYingman/arkwaifu?style=flat-square)

Arkwaifu is a website which, arranges and provides almost all picture assets extracted from Arknights (the game).

Arkwaifu also enlarges (4x) the picture assets with super-resolution neural networks, [Real-ESRGAN](https://github.com/xinntao/Real-ESRGAN)
and [Real-CUGAN](https://github.com/bilibili/ailab/tree/main/Real-CUGAN).

Currently, only assets, that appear in the "in-game stories" (what we call it: AVG), are available, including:

- "Images". They are the exquisite arts that appear when some special events happened.
- "Backgrounds". They are the backgrounds that always appear on the bottom layer, during dialogue.
- "Characters". They are the pictures of characters that act in the play, some of them are with different emotions.

I also plan to include the artwork of characters in-game.
However, by considering some technical difficulties, especially some of the artworks are animated, the plan is delayed.

This project is now available online!
Check it at [arkwaifu.cc](https://arkwaifu.cc/) (main site), or [cn.arkwaifu.cc](https://cn.arkwaifu.cc/) (CN mirror)! 🎉

# Features

- Assets are kept up-to-date automatically whenever the game is updated.
- Assets are enlarged with super-resolution neural networks ([Real-ESRGAN](https://github.com/xinntao/Real-ESRGAN)
  , [Real-CUGAN](https://github.com/bilibili/ailab/tree/main/Real-CUGAN)).
- Assets are classfied into handy categories (AVG, non-AVG; main stories, activities; etc.).

# TODOs V1 #

- [x] Beautify the frontend. (Done<del>, probably uglified</del>)
- [x] List assets which aren't included in AVGs. e.g., assets appeared in mode *Integrated Strategies*.
- [x] Use cache to speed up website (backend).
- [x] Use cache to speed up website (frontend).
- [x] Pull only differences every update loop.
- [x] I18N 🌏! Add Chinese support.
- [x] Extract gamedata directly from resources.
- [x] Provide API to manually update resources.
- [x] Assets image super-resolution. (Real-ESGRAN or Real-CUGAN)
- [x] Rewrite controller interface.
- [ ] Optimize the frontend UX.
- [ ] Support searching assets.
- [ ] Provide a choice to switch the website language.
- [ ] Detect the user's language, and automatically switch the website language to it.
- [ ] //...

# Design

For the design documentation, view [here (DESIGN.md)](DESIGN.md).

# Acknowledgements

Thanks to my friend [Galvin Gao](https://github.com/GalvinGao)!
He helped me a lot in the front-end development and choosing frameworks. I really appreciate the "getting hand dirty" methodology very much.

Thanks to my friend [Martin Wang](https://github.com/martinwang2002)!
He helped me in extracting the gamedata assets, and also in some details of automatically updating the assets from the game.

Thanks to my friend Lily! She drew the fascinating [Phantom logo](assets/arkwaifu_phantom.png) of this project.

Thanks to [Penguin Statistics](https://penguin-stats.io/)!
The prototype of this project referenced and is inspired by Penguin Statistics' [backend v3](https://github.com/penguin-statistics/backend-next).

# License

The source code of this project is licensed under the [MIT License](LICENSE).

The assets of this project are licensed under
[Attribution-NonCommercial 4.0 International (CC BY-NC 4.0)](https://creativecommons.org/licenses/by-nc/4.0/).

This project utilizes resources and other works from the game Arknights. The copyright of such works belongs to the
provider of the game, 上海鹰角网络科技有限公司 (Shanghai Hypergryph Network Technology Co., Ltd).

Some initial template source code of this project is inspired by
the [backend v3](https://github.com/penguin-statistics/backend-next) of [Penguin Statistics](https://penguin-stats.io/),
which is licensed under the [MIT License](https://github.com/penguin-statistics/backend-next/blob/dev/LICENSE).

# [x] TODOs V0 #

- [x] Follow the package style guideline. While I'm currently working on making this app running, therefore this project
  isn't following the package style
  guideline (https://github.com/danceyoung/paper-code/blob/master/package-style-guideline/packagestyleguideline.md)
- [x] Complete the *updateloop*. The updateloop updates the resources and gamedata continuously; it ensures the data is
  always up-to-date.
- [x] Complete the AVG part of ArkWaifu. In brief, the AVG part handles the requests related to gamedata (i.e., anything
  except for image resources).
- [x] Complete the Resource part of ArkWaifu. The resource part handles the requests of static resources, like image
  files and background files.
- [x] Complete the frontend of ArkWaifu. There should be a sidebar with categories: AVGs or ALL. The AVGs shows all AVG
  groups, and there are AVG stories under the AVG groups. The frontend shows all AVG resources under the user chosen
  group or story. Under the ALL category, the frontend simply shows all existing resources.
- [x] Dockerize ArkWaifu with CI (GitHub actions).
- [x] Make ArkWaifu go live!