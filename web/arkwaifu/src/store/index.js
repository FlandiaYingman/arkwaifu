import Vue from 'vue'
import Vuex from 'vuex'
import API_URL from '@/api'
import _ from 'lodash'

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    groups: [],
    groupsMap: {},
    groupsTypeMap: {},
    stories: null,
    storiesMap: {},
    storiesAssets: [],
    storiesAssetsMap: {},

    assets: [],
    assetsMap: {},
    assetsKindMap: {}
  },
  getters: {
    groupByID: (state) => (id) => state.groupsMap[id],
    storyByID: (state) => (id) => state.storiesMap[id],
    images: ({ assetsKindMap }) => assetsKindMap.images,
    backgrounds: ({ assetsKindMap }) => assetsKindMap.backgrounds
  },
  actions: {
    async updateAll ({ dispatch }) {
      dispatch('updateGroups')
      dispatch('updateStories')
      dispatch('updateAssets')
    },
    async updateGroups ({ state }) {
      return fetch(`${API_URL}/api/v0/avg/groups`)
        .then((resp) => resp.json())
        .then((groups) => {
          state.groups = groups.map(el => Object.freeze(el))
          state.groupsMap = _.keyBy(groups, el => el.id)
          state.groupsTypeMap = _.groupBy(groups, el => el.type)
        })
    },
    async updateStories ({ state }) {
      return fetch(`${API_URL}/api/v0/avg/stories`)
        .then((resp) => resp.json())
        .then((stories) => {
          state.stories = stories.map(el => Object.freeze(el))
          state.storiesMap = _.keyBy(stories, el => el.id)
          state.storiesAssets = stories.flatMap(story => {
            return story.assets.map(asset => ({
              asset: asset,
              storyId: story.id
            }))
          })
          state.storiesAssetsMap = _.keyBy(state.storiesAssets, el => el.asset.name)
        })
    },
    async updateAssets ({ state }) {
      return fetch(`${API_URL}/api/v0/asset/assets`)
        .then(resp => resp.json())
        .then(assets => {
          state.assets = assets.map(el => Object.freeze(el))
          state.assetsMap = _.keyBy(assets, el => el.name)
          state.assetsKindMap = _.groupBy(assets, el => el.kind)
        })
    }
  }
})

export default store
