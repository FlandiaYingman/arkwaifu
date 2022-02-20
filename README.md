# Ark Waifu (arkwaifu)

Automatically get art resources from Arknights with super-resolution.

This project is currently under development.

# TODOs V0 #

- [x] Follow the package style guideline. While I'm currently working on making this app running, therefore this project
  isn't following the package style
  guideline (https://github.com/danceyoung/paper-code/blob/master/package-style-guideline/packagestyleguideline.md)
- [x] Complete the *updateloop*. The updateloop updates the resources and gamedata continuously; it ensures the data is
  always up-to-date.
- [x] Complete the AVG part of Ark Waifu. In brief, the AVG part handles the requests related to gamedata (i.e.,
  anything except for image resources).
- [x] Complete the Resource part of Ark Waifu. The resource part handles the requests of static resources, like image
  files and background files.
- [ ] Complete the frontend of Ark Waifu. There should be a sidebar with categories: AVGs or ALL. The AVGs shows all AVG
  groups, and there are AVG stories under the AVG groups. The frontend shows all AVG resources under the user chosen
  group or story. Under the ALL category, the frontend simply shows all existing resources.
- [ ] Make Ark Waifu go live!
- [ ] Advertise Ark Waifu on Bilibili or somewhere...

# TODOs V1 #

- [ ] Create a statistic module to show how many times the website is requested.
- [ ] //...

# App Design #

## Update Loop ##

1. Check whether there's a newer version of Arknights periodically.
2. If there is, download and parse the updated part.
3. Save the data into database and resource into local storage.

## AVG ##

To be completed...

## Resource ##

To be completed...