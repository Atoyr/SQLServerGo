import axios from 'axios'
export const state = () => ({
  instanceName : ""
})

export const mutations = {
  updateInstanceName(state, props) {
    state.instanceName = props.instanceName;
  }
}

export const getters = {
  instanceName (state) {
    return state.instanceName;
  }
}

export const actions = {
  async fetchInstance({commit}) {
    this.$axios.$get(`http://${this.$getHost()}/api/instance`)
    .then((res) => {
      commit('updateInstanceName',{instanceName:res.serverName});
    })
  }
}

