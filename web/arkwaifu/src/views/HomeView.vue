<template>
  <v-container>
    <img
      src="/logo_phantom.png"
      alt=""
    >
    <p class="text-h3">
      {{ $t("title") }}
    </p>
    <p class="text-subtitle-1">
      {{ $t("description") }}
    </p>
    <p class="text-subtitle-1">
      {{ $t("prompt") }}
    </p>
  </v-container>
</template>

<script>
export default {
  name: 'HomeView',
  components: {}
}
</script>

<i18n>
{
  "en": {
    "title": "Welcome to Arkwaifu!",
    "description": "Arkwaifu provides the art resources get from Arknights. ",
    "prompt": "Find the resource using the navigation drawer on the left. Enjoy it!"
  },
  "zh": {
    "title": "欢迎来到 Arkwaifu！",
    "description": "Arkwaifu 提供从《明日方舟》中获取的美术资源。",
    "prompt": "用左边的导航栏开始寻找资源。玩的开心！"
  }
}
</i18n>
