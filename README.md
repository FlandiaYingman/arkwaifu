# ArkWaifu (arkwaifu)

![](https://img.shields.io/github/license/FlandiaYingman/arkwaifu?style=flat-square)
![](https://img.shields.io/github/last-commit/FlandiaYingman/arkwaifu?style=flat-square)

A website providing AVG art resources of Arknights.

This project is currently under development.

# TODOs V0 #

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
- [ ] Dockerize ArkWaifu.
- [ ] Make ArkWaifu go live!
- [ ] Advertise ArkWaifu on Bilibili or somewhere...

# TODOs V1 #

- [ ] Create a statistic module to show how many times the website is requested.
- [ ] //...

# App Design #

## Frontend

In left there's a sidebar like [Wikipedia](https://wikipedia.org/) or any IDEs. And 3 categories inside the sidebar - "
Home", "AVGs" and "All". On top of the sidebar there's a search bar to search any *groups* or *stories*.

For example, the sidebar should look like this:

```
Home
AVGs
├── Mainline
│   ├── ...
│   ├── 怒号光明
│   └── 风暴瞭望
├── Acivities
│   ├── 阴云火花
│   ├── 将进酒
│   └── ...
└── Operator Record
    ├── 学者之心
    ├── 火山
    ├── 特大号烤饼
    └── ...
All
├── Images
└── Backgrounds
```

And see the following for explanation.

### Home

TODO...

### AVGs

A *group* is a group of story, such as "将进酒" or "怒号光明". There are 3 types of AVG groups - "Mainline" (well, probably "
主线" or "Main Storyline", however Arknights call it "Mainline"), "Activity" ("活动" or "Event") and "Operator Record" ("
干员密录"). A *story* is just literally a story, i.e., an AVG. The stories in different operations are counted individually.
And the stories before and after an operation are counted individually also.

The "AVGs" category is for exploring the AVG resources categorized by AVG *groups* and *stories*. Therefore, under the "
AVGs" category there are the different group type ("Mainline" etc.). Under group types there are *groups* and under
groups there are *stories*.

The group types and groups can't be chosen. The user can only expand or fold them.

Choosing a story would show all image and background resources it uses.

I haven't figured out how does it show image and background resources... Just simply show thumbnails list now.

### All

Same to AVGs. Just simply show thumbnails...

## Backend

### Update Loop ###

1. Check whether there's a newer version of Arknights periodically.
2. If there is, download and parse the updated part.
3. Save the data into database and resource into local storage.

### AVG ###

TODO

#### Group

TODO

#### Story

TODO

### Resource ###

There are two kinds of resource - image and background.

#### Image

TODO

#### Background

TODO
