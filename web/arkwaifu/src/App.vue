<template>
  <main-layout>
    <keep-alive>
      <router-view :key="$route.fullPath" />
    </keep-alive>
  </main-layout>
</template>

<script>

import MainLayout from '@/layouts/MainLayout'

export default {
  name: 'App',
  components: { MainLayout },
  data: () => ({
    title: '',
    drawer: null
  }),
  watch: {
    $route: {
      handler (to) {
        this.title = to.meta.title
      },
      immediate: true
    }
  }
}
</script>
