<template>
  <div>
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
    <p class="text-h4">
      {{ $t("donate.title") }}
    </p>
    <p class="text-subtitle-1">
      {{ $t("donate.text") }}
    </p>
    <img src="/wechat_pay_denotation.png">
    <p class="text-sm-caption">
      <del> {{ $t("donate.easter_egg") }}</del>
    </p>
  </div>
</template>

<script>

export default {
  name: 'HomeView'
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
    "prompt": "用左边的导航栏开始寻找资源。玩的开心！",
    "donate.title": "捐助",
    "donate.text": "作为资源站，即便租用的是国外服务器，服务器费用也是不可忽视的。目前，服务器每月的费用是 10 USD（折合人民币约 70 元）。如果各位用户能捐助一些服务器费用，将不胜感激！另外，捐助请注明来源，以帮助更好地计算收支。",
    "donate.easter_egg": "作为对捐助者的福利……请访问站长 GitHub 某仓库"
  }
}
</i18n>
