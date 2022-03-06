<!--suppress ALL -->
<img src="assets/arkwaifu_phantom@0.25x.png" alt="logo" align="right" height="224" width="224"/>

# App Design #

## Backend

### Concepts

- **Assets**: graphics or audio assets of the game, such as background pictures or character illustrations; strictly
  speaking, gamedata is a part of the assets, but generally, in the ArkWaifu context, it's not.
- **Resources**: Same as *asset*; used before ArkWaifu version 0.1.0.
- **Gamedata; Data**: the gamedata (text) of the game, such as story scripts or operator tables.
- **AVG**: the sum of all stories in the game; one can say AVG assets and AVG gamedata, but seldom story assets or story
  gamedata.
- **Story**: an individual story in the game; e.g., the story before 1-7, the story after 1-7, a personal story in a
  Vignette.
- **Story Group; Group**: a logical group for stories; for the Main Themes, they are the Episodes like *Burning Run* (
  Episode 04) etc.; for the Events, they are the Events themselves, like *Under Tides* or *Children of Ursus* etc.; for
  Operator Records, they are also the Operator Records themselves (note: 2 records of 1 operator are counted as 2
  groups; 2 stories of 1 record are counted as 1 group).

#### Assets

- **Image**: often refer to the `image` category in AVG assets.
- **Background**: often refer to the `background` category in AVG assets.

### Update Loop (updateloop) ###

The duty of updateloop module is to keep assets and gamedata up-to-date.

Whenever the game updates, updateloop pull the latest assets and gamedata, and push then into the local storage (for
assets) and database (for gamedata).

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