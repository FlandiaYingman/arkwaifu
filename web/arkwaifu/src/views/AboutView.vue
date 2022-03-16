<template>
  <v-container>
    <p
      class="text-h3"
    >
      {{ $t('title') }}
    </p>
    <p
      class="text-body-1"
    >
      {{ $t("p1") }}
    </p>
    <i18n
      class="text-body-1"
      path="p2"
      tag="p"
    >
      <template #repo>
        <a
          href="https://github.com/FlandiaYingman/arkwaifu"
          target="_blank"
        >https://github.com/FlandiaYingman/arkwaifu</a>
      </template>
    </i18n>
  </v-container>
</template>

<script>

export default {

}
</script>

<i18n>
{
  "en": {
    "title": "About Arkwaifu",
    "p1": "Hi, I'm Flandia Yingman, Arkwaifu's author. This project (Arkwaifu) has been founded for only a few days so there's probably nothing to say in \"about\" :D.",
    "p2": "However, if you are still interested about Arkwaifu, please go to the Github repo ({repo}). Any kinds of issue or pull requests are welcome. Thanks!"
  },
  "zh": {
    "title": "关于 Arkwaifu",
    "p1": "你好！我是 Arkwaifu 的作者孔明亲 (@FlandiaYingman)。Arkwaifu 才刚被创造几天，所以没啥可说的！XD",
    "p2": "但是，如果你真的很感兴趣，请来 GitHub ({repo}) 上提 Issue 或者 Pull Request 吧！"
  }
}
</i18n>
