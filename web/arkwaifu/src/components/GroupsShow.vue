<template>
  <div>
    <group-show
      v-for="(group, i) in groups"
      :key="i"
      :group-id="group.id"
    />
    <fab-button
      v-model="descending"
      icon-on="mdi-sort-calendar-descending"
      icon-off="mdi-sort-calendar-ascending"
      caption-on="Sort: Descending"
      caption-off="Sort: Ascending"
    />
  </div>
</template>

<script>
import GroupShow from '@/components/GroupShow'
import FabButton from '@/components/FabButton'
import _ from 'lodash'

export default {
  name: 'GroupsShow',
  components: { GroupShow, FabButton },
  props: {
    type: String()
  },
  data () {
    return {
      descending: true
    }
  },
  computed: {
    groups () {
      let groups = this.$store.state.groupsTypeMap[this.type]
      if (groups) {
        if (this.descending) {
          groups = _.reverse(groups)
        }
        return groups
      }
      return null
    }
  }
}
</script>
