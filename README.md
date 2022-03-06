<!--suppress ALL -->
<img src="assets/arkwaifu_phantom@0.25x.png" alt="logo" align="right" height="224" width="224"/>

# ArkWaifu (arkwaifu)

[![](https://pkg.go.dev/badge/github.com/flandiayingman/arkwaifu.svg)](https://pkg.go.dev/github.com/flandiayingman/arkwaifu)
![](https://img.shields.io/github/license/FlandiaYingman/arkwaifu?style=flat-square)
![](https://img.shields.io/github/last-commit/FlandiaYingman/arkwaifu?style=flat-square)

A website providing AVG art resources of Arknights.

This project is currently under development.

# TODOs V1 #

- [ ] Advertise ArkWaifu on Bilibili or somewhere...
- [ ] Create a statistic module to show how many times the website is requested.
- [ ] //...

# Design

For the design documentation, view at [here (DESIGN.md)](DESIGN.md).

# Acknowledgements

Thanks to Galvin Gao! He helped me a lot with the front-end development and in choosing frameworks. I really appreciate
the "getting hand dirty" methodology very much.

Thanks to Penguin Statistics! The prototype of this project has referenced and is inspired by Penguin
Statistics' [backend v3](https://github.com/penguin-statistics/backend-next).

Thanks to my friend Lily! She drew the fascinating Phantom logo of this project.

# License

The source code of this project is licensed under the [MIT License](LICENSE).

The assets of this project is licensed
under [Attribution-NonCommercial 4.0 International (CC BY-NC 4.0)](https://creativecommons.org/licenses/by-nc/4.0/).

This project utilizes resources and other works from the game Arknights. The copyright of such works belongs to the
provider of the game, Shanghai Hypergryph Network Technology Co., Ltd.

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