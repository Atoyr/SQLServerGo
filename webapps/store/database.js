import axios from 'axios'

export const state = () => ({
  serverProperty: {
    machineName: "",
	  instanceName: "",
	  serverName: "",
	  productVersion: "",
	  productMajorVersion: "",
	  version: "",
	  edition: "",
	  productLevel: ""
  },
  databaseFiles: [],
  fileInputIO: [],
  fileOutputIO: [],
})

export const mutations = {
  updateServerProperty(state, props) {
    state.serverProperty.machineName = props.serverProperty.machineName;
    state.serverProperty.instanceName = props.serverProperty.instanceName;
    state.serverProperty.serverName = props.serverProperty.serverName;
    state.serverProperty.productVersion = props.serverProperty.productVersion;
    state.serverProperty.productMajorVersion = props.serverProperty.productMajorVersion;
    state.serverProperty.version = props.serverProperty.version;
    state.serverProperty.edition = props.serverProperty.edition;
    state.serverProperty.productLevel = props.serverProperty.productLevel;
  }
}

export const getters = {
  InstanceName (state) {
    return state.serverProperty.instanceName;
  }
}

export const actions = {
  async fetchServerProperty ({commit}) {
    await axios.get(`http://${this.$getHost()}/api/instance`)
    .then((res) => {
      commit('updateServerProperty',{serverProperty:res.data});
    })
  }
}
