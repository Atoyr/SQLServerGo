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
  instance: {},
  cpu: {
    id: "",
    systemIdle: "",
    sqlProcessUtilization: "",
    otherProcessUtilization: "",
    timestamp:""
  },
  serverStatus: {
    startTime: new Date()
  }
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
  },
  updateInstance(state, props) {
    let datetime = new Date(props.data.datetime);
    // create database struct in instance
    props.data.databases.forEach((d) => {
      if (d.database_name in state.instance){
      } else{
        state.instance[d.database_name] = {files:[],read:[],write:[]}
        for(const f of d.files){
          state.instance[d.database_name].files.push(f.file_name)
        }
      }
    });

    // add X key
    for(let key in state.instance) {
      state.instance[key].read.push([datetime]);
      state.instance[key].write.push([datetime]);
    }

    // add value exchange KiB 
    for(const d of props.data.databases){
      let index = state.instance[d.database_name].read.length - 1
      for(const f of d.files){
        state.instance[d.database_name].read[index].push(f.read_bytes_per_sec / 1024)
        state.instance[d.database_name].write[index].push(f.write_bytes_per_sec / 1024)
      }
    }
    
    // apply data max length 
    for(let key in state.instance) {
      if (state.instance[key].read.length > 30) {
        state.instance[key].read.shift();
      }
      if (state.instance[key].write.length > 30) {
        state.instance[key].write.shift();
      }
    }
  },
  updateServerStatus(state, props) {
    state.serverStatus.startTime = new Date(props.serverStatus.instance_start_time);
  }
}

export const getters = {
  Instance (state) {
    return state.instance;
  },
  InstanceName (state) {
    return state.serverProperty.instanceName;
  },
  ServerName (state) {
    return state.serverProperty.serverName;
  },
  StartTime (state) {
    return state.serverStatus.startTime;
  }
}

export const actions = {
  async fetchServerProperty ({commit}) {
    await axios.get(`http://${this.$getHost()}/api/instance`)
    .then((res) => {
      console.log(res.data)
      commit('updateServerProperty',{serverProperty:res.data});
    })
  },
  async fetchDatabaseFiles({commit}) {
    await axios.get(`http://${this.$getHost()}/api/databaseFiles`)
    .then((res) => {
      console.log(res.data)
    })
  },
  async fetchServerStatus({commit}) {
    await axios.get(`http://${this.$getHost()}/api/serverStatus`)
    .then((res) => {
      commit('updateServerStatus',{serverStatus:res.data})
      console.log(res.data)
    })
  }
}
