<!--suppress ALL -->
<img src="assets/arkwaifu_phantom@0.25x.png" alt="logo" align="right" height="224" width="224"/>

# ArkWaifu (arkwaifu)

[![](https://pkg.go.dev/badge/github.com/flandiayingman/arkwaifu.svg)](https://pkg.go.dev/github.com/flandiayingman/arkwaifu)
![](https://img.shields.io/github/license/FlandiaYingman/arkwaifu?style=flat-square)
![](https://img.shields.io/github/last-commit/FlandiaYingman/arkwaifu?style=flat-square)

Arkwaifu is a website which, arranges and provides almost all picture assets extracted from Arknights (the game).

Arkwaifu also enlarges (4x) the picture assets with super-resolution neural
networks, [Real-ESRGAN](https://github.com/xinntao/Real-ESRGAN)
and [Real-CUGAN](https://github.com/bilibili/ailab/tree/main/Real-CUGAN).

Currently, only assets, that appear in the "in-game stories" (what we call: AVG), are available, including:

- "Images". They are the exquisite arts that appear when some special events happened.
- "Backgrounds". They are the backgrounds that always appear on the bottom layer, during dialogue.
- "Characters". They are the pictures of characters that act in the play, some of them are with different emotions.

I also plan to include the artwork of characters in-game.
However, by considering some technical difficulties, especially there are some animated artworks, the plan is delayed.

This project is now available online!
Check it at [arkwaifu.cc](https://arkwaifu.cc/) (main site), or [cn.arkwaifu.cc](https://cn.arkwaifu.cc/) (CN mirror)!
🎉

# Features

- Assets are kept up-to-date automatically whenever the game is updated.
- Assets are enlarged with super-resolution neural networks ([Real-ESRGAN](https://github.com/xinntao/Real-ESRGAN)
  , [Real-CUGAN](https://github.com/bilibili/ailab/tree/main/Real-CUGAN)).
- Assets are classified into handy categories (AVG, non-AVG; main stories, activities; etc.).
- Containerized backend (microservices) and frontend.
- The backend follows microservice architecture and is split into two services: 'service', which serves the data, and
  'updateloop' which updates the data automatically.

# V1 Roadmap #

- [x] Beautify the frontend.
- [x] List assets which aren't included in AVGs.
- [x] Use cache to speed up the website (backend).
- [x] Use cache to speed up the website (frontend).
- [x] Pull only differences in every update loop.
- [x] I18N 🌏! Add Chinese support.
- [x] Extract gamedata directly from the game resources.
- [x] Create API to manually update the resources.
- [x] Support enlarging picture assets with neural networks.
- [x] Redesign & Rewrite the controller interface.
- [ ] Redesign the frontend UX.
- [ ] Support assets searching.
- [ ] Provide a choice to switch the website language.
- [ ] Automatically switch the website language to user's language automatically.
- [ ] //...

# Design

For the design documentation, see [here (DESIGN.md)](DESIGN.md).

# Acknowledgements

Thanks to my friend [Galvin Gao](https://github.com/GalvinGao)!
He helped me a lot in the front-end development and choosing frameworks. I really appreciate the "getting hands dirty"
methodology very much.

Thanks to my friend [Martin Wang](https://github.com/martinwang2002)!
He helped me in extracting the gamedata assets, and also in some details of automatically updating the assets from the
game.

Thanks to my friend Lily! She drew the fascinating [Phantom logo](assets/arkwaifu_phantom.png) of this project.

Thanks to [Penguin Statistics](https://penguin-stats.io/)!
The prototype of this project referenced and is inspired by Penguin
Statistics' [backend v3](https://github.com/penguin-statistics/backend-next).

Thanks to [xinntao](https://github.com/xinntao), [nihui](https://github.com/nihui), and the other contributors
of [Real-ESRGAN](https://github.com/xinntao/Real-ESRGAN)
and [Real-CUGAN](https://github.com/bilibili/ailab/tree/main/Real-CUGAN)! They created the neural networks this project
utilizes for enlarging assets.

# License

The source code of this project is licensed under the [MIT License](LICENSE).

The assets of this project are licensed under
[Attribution-NonCommercial 4.0 International (CC BY-NC 4.0)](https://creativecommons.org/licenses/by-nc/4.0/).

This project utilizes resources and other works from the game Arknights. The copyright of such works belongs to the
provider of the game, 上海鹰角网络科技有限公司 (Shanghai Hypergryph Network Technology Co., Ltd).

This project utilizes [Real-ESRGAN](https://github.com/xinntao/Real-ESRGAN)
and [Real-ESRGAN-ncnn-vulkan](https://github.com/xinntao/Real-ESRGAN-ncnn-vulkan), which are respectively licensed under
the BSD-3-Clause license and the MIT License.

This project utilizes [Real-CUGAN](https://github.com/bilibili/ailab/tree/main/Real-CUGAN)
and [Real-CUGAN-ncnn-vulkan](https://github.com/nihui/realcugan-ncnn-vulkan), which are both licensed under the MIT
License.

Some initial template source code of this project is inspired by and copied from
the [backend v3](https://github.com/penguin-statistics/backend-next) of [Penguin Statistics](https://penguin-stats.io/),
which is licensed under the [MIT License](https://github.com/penguin-statistics/backend-next/blob/dev/LICENSE).
