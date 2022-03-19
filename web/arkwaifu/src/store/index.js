import Vue from 'vue'
import Vuex from 'vuex'
import API_URL from '@/api'
import _ from 'lodash'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    avg: {
      state: {
        groups: null,
        stories: null,
        groupsTypeMap: {}
      },
      getters: {
        groupByID: (state) => (id) => state.groups.find((el) => el.id === id),
        storyByID: (state) => (id) => state.stories.find((el) => el.id === id)
      },
      actions: {
        async updateAll ({ dispatch }) {
          dispatch('updateGroups')
          dispatch('updateStories')
        },
        async updateGroups (store) {
          return fetch(`${API_URL}/api/v0/groups`)
            .then((resp) => resp.json())
            .then((groups) => {
              store.state.groups = groups.map(el => Object.freeze(el))
              store.state.groupsTypeMap = _.groupBy(groups, el => el.actType)
            })
        },
        async updateStories ({ state }) {
          return fetch(`${API_URL}/api/v0/stories`)
            .then((resp) => resp.json())
            .then((stories) => (state.stories = stories.map(el => Object.freeze(el))))
        }
      }
    },
    assets: {
      state: {
        assets: []
      },
      getters: {
        images: state => {
          return state.assets.filter(el => el.kind === 'images')
        },
        backgrounds: state => {
          return state.assets.filter(el => el.kind === 'backgrounds')
        }
      },
      actions: {
        async updateAll ({ dispatch }) {
          dispatch('updateAssets')
        },
        async updateAssets ({ state }) {
          return fetch(`${API_URL}/api/v0/assets/img`)
            .then((resp) => resp.json())
            .then((assets) => (state.assets = assets.map(el => Object.freeze(el))))
        }
      }
    }
  }
})

export default store
