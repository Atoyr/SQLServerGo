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
    let datetime = new Date(props.data[0].datetime);
    let fileCount = 0;
    let beforeDatabaseName = ""
    // create database struct in instance
    for(const f of props.data){
      if (beforeDatabaseName != f.database_name) {
        fileCount = 0;
        beforeDatabaseName = f.database_name;
      }
      if (f.database_name in state.instance){
        if (state.instance[f.database_name].files.length <= fileCount ) {
          state.instance[f.database_name].files.push(f.file_name)
        }
      }else {
        state.instance[f.database_name] = {files:[f.file_name],read:[],write:[]}
      }
      fileCount++;
    }

    // add X key
    for(let key in state.instance) {
      state.instance[key].read.push([datetime]);
      state.instance[key].write.push([datetime]);
    }
    // add value exchange KiB 
    for(const f of props.data){
      let index = state.instance[f.database_name].read.length - 1
      state.instance[f.database_name].read[index].push(f.read_bytes_per_sec / 1024)
      state.instance[f.database_name].write[index].push(f.write_bytes_per_sec / 1024)
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
