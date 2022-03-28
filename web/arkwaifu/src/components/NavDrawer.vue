<template>
  <v-navigation-drawer
    width="320px"
    permanant
    app
    :value="value"
    @input="$emit('input', $event )"
  >
    <v-card
      v-ripple
      class="ma-4 pa-2"
      elevation="0"
      rounded="xl"
      color="accent lighten-5"
      style=""
    >
      <div class="mr-5">
        <v-img
          src="/logo_phantom.png"
          width="192px"
          height="192px"
          class="mx-auto"
          style="filter: drop-shadow(-15px 15px 15px rgba(0,0,0,0.15)); user-select: none;"
        />
      </div>
      <v-card-title class="justify-center text-caption">
        <span
          class="primary--text"
          style="user-select: none;"
        >Arkwaifu - Phantom</span>
      </v-card-title>
    </v-card>

    <v-list nav>
      <v-subheader>
        {{ $t("header") }}
      </v-subheader>
      <v-list-item-group
        color="primary"
      >
        <v-list-item
          link
          to="/"
        >
          <v-list-item-icon>
            <v-icon>{{ mdiHome }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              {{ $t("home.title") }}
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-group
          link
          :prepend-icon="mdiTelevision"
        >
          <template #activator>
            <v-list-item-content>
              <v-list-item-title>
                {{ $t("avg.title") }}
              </v-list-item-title>
              <v-list-item-subtitle>
                {{ $t("avg.subtitle") }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </template>
          <v-list-item to="/avgs/main_themes">
            <v-list-item-icon>
              <v-icon>{{ mdiCube }}</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>
                {{ $t("avg.mainThemes.title") }}
              </v-list-item-title>
              <v-list-item-subtitle>
                {{ $t("avg.mainThemes.subtitle") }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-list-item to="/avgs/major_events">
            <v-list-item-icon>
              <v-icon>{{ mdiMusicClefTreble }}</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>
                {{ $t("avg.majorEvents.title") }}
              </v-list-item-title>
              <v-list-item-subtitle>
                {{ $t("avg.majorEvents.subtitle") }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-list-item to="/avgs/vignettes">
            <v-list-item-icon>
              <v-icon>{{ mdiMusicClefBass }}</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>
                {{ $t("avg.vignettes.title") }}
              </v-list-item-title>
              <v-list-item-subtitle>
                {{ $t("avg.vignettes.subtitle") }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-list-item to="/avgs/others">
            <v-list-item-icon>
              <v-icon>{{ mdiMusicRestQuarter }}</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>
                {{ $t("avg.others.title") }}
              </v-list-item-title>
              <v-list-item-subtitle>
                {{ $t("avg.others.subtitle") }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list-group>
        <v-list-item
          link
          to="/non_avgs"
        >
          <v-list-item-icon>
            <v-icon>{{ mdiTelevisionOff }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              {{ $t("nonAvg.title") }}
            </v-list-item-title>
            <v-list-item-subtitle>
              {{ $t("nonAvg.subtitle") }}
            </v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
        <v-list-item
          link
          to="/all"
        >
          <v-list-item-icon>
            <v-icon>{{ mdiAllInclusive }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              {{ $t("all.title") }}
            </v-list-item-title>
            <v-list-item-subtitle>
              {{ $t("all.subtitle") }}
            </v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
        <v-list-item
          link
          to="/about"
        >
          <v-list-item-icon>
            <v-icon>{{ mdiInformation }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              {{ $t("about.title") }}
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import {
  mdiAllInclusive,
  mdiCube,
  mdiHome,
  mdiInformation,
  mdiMusicClefBass,
  mdiMusicClefTreble,
  mdiMusicRestQuarter,
  mdiTelevision,
  mdiTelevisionOff
} from '@mdi/js'

export default {
  name: 'NavigationDrawer',
  props: {
    value: Boolean
  },
  data: () => ({
    selected: {},
    mdiHome: mdiHome,
    mdiCube: mdiCube,
    mdiMusicClefTreble: mdiMusicClefTreble,
    mdiMusicClefBass: mdiMusicClefBass,
    mdiMusicRestQuarter: mdiMusicRestQuarter,
    mdiTelevision: mdiTelevision,
    mdiTelevisionOff: mdiTelevisionOff,
    mdiAllInclusive: mdiAllInclusive,
    mdiInformation: mdiInformation
  })
}
</script>

<i18n>
{
  "en": {
    "header": "NAVIGATION",
    "home": {
      "title": "Home"
    },
    "avg": {
      "title": "AVG",
      "subtitle": "Sort by AVG",
      "mainThemes": {
        "title": "Main Themes",
        "subtitle": "Main Stories"
      },
      "majorEvents": {
        "title": "Major Events",
        "subtitle": "Intermezzi & Side Stories"
      },
      "vignettes": {
        "title": "Vignettes",
        "subtitle": "Story Collections"
      },
      "others": {
        "title": "Others",
        "subtitle": "Operator Records"
      }
    },
    "nonAvg": {
      "title": "Non-AVG",
      "subtitle": "Besides AVG"
    },
    "all": {
      "title": "All Assets",
      "subtitle": "May take a long time! "
    },
    "about": {
      "title": "About"
    }
  },
  "zh": {
    "header": "导航栏",
    "home": {
      "title": "首页"
    },
    "avg": {
      "title": "剧情",
      "subtitle": "剧情索引资源",
      "mainThemes": {
        "title": "主题曲",
        "subtitle": "主线剧情"
      },
      "majorEvents": {
        "title": "活动",
        "subtitle": "插曲 & 别传"
      },
      "vignettes": {
        "title": "故事集",
        "subtitle": "故事集 & 微型故事集"
      },
      "others": {
        "title": "其它",
        "subtitle": "干员密录"
      }
    },
    "nonAvg": {
      "title": "非剧情",
      "subtitle": "剧情以外的资源"
    },
    "all": {
      "title": "全部资源",
      "subtitle": "小心！加载时间较长"
    },
    "about": {
      "title": "关于"
    }
  }
}
</i18n>
