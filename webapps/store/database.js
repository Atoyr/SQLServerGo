import axios from 'axios'
import { w3cwebsocket } from 'websocket';
const W3CWebSocket = w3cwebsocket

const fileIOWebsocketPlugin = (store) => {
    var ws = new W3CWebSocket(`ws://${this.$getHost()}/ws/fileio`)
    ws.onmessage = (e) => {
      if (typeof e.data === 'string') {
        let data = JSON.parse(event.data);
        if (data.length > 0){
          this.labels.push( Date.parse(data[0].datetime));
        }
        for(const f of data){
          let id = (f.id - 1) * 2;

          this.datasets[id].data.push({
            x: Date.parse(f.datetime),
            y: f.read_bytes_per_sec,
          });
          this.datasets[id + 1].data.push({
            x: Date.parse(f.datetime),
            y: f.read_bytes_per_sec,
          });

          if (this.datasets[id].data.length > 120) {
            this.datasets[id].data.pop();
          }
          if (this.datasets[id + 1].data.length > 120) {
            this.datasets[id + 1].data.pop();
          }
        }
        if (this.labels.length > 120) {
          this.labels.pop();
        }
      }
    }
}

const cpuWebsocketPlugin = (store) => {
    var ws = new W3CWebSocket(`ws://${this.$getHost()}/ws/cpu`)
    ws.onmessage = (e) => {
      if (typeof e.data === 'string'){
        let data = JSON.parse(event.data);
        console.log(`cpu plugin ${data}`)
      }
    }
}

export const plugins = [
  fileIOWebsocketPlugin,
  cpuWebsocketPlugin
]


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
  updateFileIO({index,time,input,output}){
  },
  updateServerStatus(state, props) {
    state.serverStatus.startTime = new Date(props.serverStatus.instance_start_time);
  }
}

export const getters = {
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
      //commit('updateServerProperty',{serverProperty:res.data});
    })
  },
  async fetchServerStatus({commit}) {
    await axios.get(`http://${this.$getHost()}/api/serverStatus`)
    .then((res) => {
      commit('updateServerStatus',{serverStatus:res.data})
      // instance_start_time
      console.log(res.data)
      //commit('updateServerProperty',{serverProperty:res.data});
    })
  }
}
