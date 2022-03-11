import Vue from "vue";
import Vuex from "vuex";
import { API_URL } from "@/api";

Vue.use(Vuex);

const store = new Vuex.Store({
  modules: {
    avg: {
      state: {
        groups: [],
        stories: [],
      },
      getters: {
        groupByID: (state) => (id) => {
          return state.groups.find((el) => el.id == id);
        },
        storyByID: (state) => (id) => {
          return state.stories.find((el) => el.id == id);
        },
      },
      mutations: {
        setGroups(state, payload) {
          state.groups = payload;
        },
        setStories(state, payload) {
          state.stories = payload;
        },
      },
      actions: {
        async updateAll({ dispatch }) {
          dispatch("updateGroups");
          dispatch("updateStories");
        },
        async updateGroups({ commit }) {
          return fetch(`${API_URL}/api/v0/groups`)
            .then((resp) => resp.json())
            .then((groups) => commit("setGroups", groups));
        },
        async updateStories({ commit }) {
          return fetch(`${API_URL}/api/v0/stories`)
            .then((resp) => resp.json())
            .then((stories) => commit("setStories", stories));
        },
      },
    },
  },
});

export default store;