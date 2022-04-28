<template>
  <div>
    <v-lazy
      v-for="(group, i) in groups"
      :key="i"
    >
      <group-show :group-id="group.id" />
    </v-lazy>
    <fab-button
      v-model="descending"
      :icon-on="mdiSortCalendarDescending"
      :icon-off="mdiSortCalendarAscending"
      :caption-on="$t('sortDesc')"
      :caption-off="$t('sortAsc')"
    />
  </div>
</template>

<script>
import GroupShow from '@/components/GroupShow'
import FabButton from '@/components/FabButton'
import _ from 'lodash'
import { mdiSortCalendarAscending, mdiSortCalendarDescending } from '@mdi/js'

export default {
  name: 'GroupsShow',
  components: { GroupShow, FabButton },
  props: {
    type: String()
  },
  data () {
    return {
      descending: true,
      mdiSortCalendarDescending: mdiSortCalendarDescending,
      mdiSortCalendarAscending: mdiSortCalendarAscending
    }
  },
  computed: {
    groups () {
      let groups = [...this.$store.getters.groupsByType(this.type)]
      if (this.descending) {
        groups = _.reverse(groups)
      }
      return groups
    }
  }
}
</script>

<i18n>{
  "en": {
    "sortAsc": "Sort: Ascending",
    "sortDesc": "Sort: Descending"
  },
  "zh": {
    "sortAsc": "排序: 升序",
    "sortDesc": "排序: 降序"
  }
}</i18n>
