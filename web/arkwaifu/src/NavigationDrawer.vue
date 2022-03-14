<template>
  <v-navigation-drawer
    width="288"
    permanant
    app
    :value="value"
    @input="$emit('input', $event.target.value)"
  >
    <v-img
      src="/logo_phantom.png"
      width="128"
      height="128"
      class="mx-4"
    />

    <v-list nav>
      <v-subheader>{{ $t("navigation") }}</v-subheader>

      <v-list-item
        link
        to="/"
      >
        <v-list-item-icon>
          <v-icon>mdi-home</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>{{ $t('home.title') }}</v-list-item-title>
          <router-link to="/home" />
        </v-list-item-content>
      </v-list-item>

      <v-list-group
        link
        prepend-icon="mdi-television"
      >
        <template #activator>
          <v-list-item-content>
            <v-list-item-title>{{ $t('avg.title') }}</v-list-item-title>
            <v-list-item-subtitle>{{ $t('avg.subtitle') }}</v-list-item-subtitle>
          </v-list-item-content>
        </template>
        <v-list-item to="/avgs/main_themes">
          <v-list-item-icon>
            <v-icon>mdi-cube</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{ $t('avg.mainThemes.title') }}</v-list-item-title>
            <v-list-item-subtitle>{{ $t('avg.mainThemes.subtitle') }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/avgs/major_events">
          <v-list-item-icon>
            <v-icon>mdi-music-clef-treble</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{ $t('avg.majorEvents.title') }}</v-list-item-title>
            <v-list-item-subtitle>{{ $t('avg.majorEvents.subtitle') }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/avgs/vignettes">
          <v-list-item-icon>
            <v-icon>mdi-music-clef-bass</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{ $t('avg.vignettes.title') }}</v-list-item-title>
            <v-list-item-subtitle>{{ $t('avg.vignettes.subtitle') }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/avgs/others">
          <v-list-item-icon>
            <v-icon>mdi-music-rest-quarter</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{ $t('avg.others.title') }}</v-list-item-title>
            <v-list-item-subtitle>{{ $t('avg.others.subtitle') }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-list-group>

      <v-list-item
        link
        to="/non_avgs"
      >
        <v-list-item-icon>
          <v-icon>mdi-television-off</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>{{ $t('nonAvg.title') }}</v-list-item-title>
          <v-list-item-subtitle>{{ $t('nonAvg.subtitle') }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        to="/all"
      >
        <v-list-item-icon>
          <v-icon>mdi-all-inclusive</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>{{ $t('all.title') }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        to="/about"
      >
        <v-list-item-icon>
          <v-icon>mdi-information</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>{{ $t('about.title') }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
export default {
  name: 'NavigationDrawer',
  props: {
    value: Boolean
  }
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
      "title": "All"
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
      "title": "全部"
    },
    "about": {
      "title": "关于"
    }
  }

}
</i18n>
