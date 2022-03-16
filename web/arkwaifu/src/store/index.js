import Vue from 'vue'
import Vuex from 'vuex'
import API_URL from '@/api'
import _ from 'lodash'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    avg: {
      state: {
        groups: [],
        stories: [],
        groupsTypeMap: {}
      },
      getters: {
        groupByID: (state) => (id) => state.groups.find((el) => el.id === id),
        storyByID: (state) => (id) => state.stories.find((el) => el.id === id)
      },
      mutations: {
        setGroups (state, payload) {
          state.groups = payload
        },
        setStories (state, payload) {
          state.stories = payload
        }
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
              store.state.groups = groups
              store.state.groupsTypeMap = _.groupBy(groups, el => el.actType)
            })
        },
        async updateStories ({ commit }) {
          return fetch(`${API_URL}/api/v0/stories`)
            .then((resp) => resp.json())
            .then((stories) => commit('setStories', stories))
        }
      }
    },
    assets: {
      state: {
        images: [],
        backgrounds: []
      },
      actions: {
        async updateAll ({ dispatch }) {
          dispatch('updateImages')
          dispatch('updateBackgrounds')
        },
        async updateImages (context) {
          return fetch(`${API_URL}/api/v0/resources/images`)
            .then((resp) => resp.json())
            .then((images) => {
              const { state } = context
              state.images = images
            })
        },
        async updateBackgrounds (context) {
          return fetch(`${API_URL}/api/v0/resources/backgrounds`)
            .then((resp) => resp.json())
            .then((backgrounds) => {
              const { state } = context
              state.backgrounds = backgrounds
            })
        }
      }
    }
  }
})

export default store
