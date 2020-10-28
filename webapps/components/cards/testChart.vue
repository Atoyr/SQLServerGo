<template>
  <v-card 
    color="primary"
    dark>
    <v-card-title class="title accent py-1">
     File Read bytes/sec
    </v-card-title>
    <div class="px-4">
      <GChart
        type="LineChart"
        :data="chartData"
        :options="chartOptions"
        />
    </div>
  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import { w3cwebsocket } from 'websocket';
import { GChart } from 'vue-google-charts';
const W3CWebSocket = w3cwebsocket
const datalength = 300

export default {
  components: {
    GChart
  },
  methods: {
  },
  computed: {
    chartData: function() {
      return this.header.concat(this.dtil);
    }
  },
  props: {
    size: Number
  },
  mounted() {
    fetch(`http://${this.$getHost()}/api/databaseFiles`)
      .then(response => {
        return response.json();
      })
      .then(res => {
        let header = ["datetime"]
        for(const f of res){
          header.push(f.file_name)
        }
        this.header = [header]
      })

    var ws = new W3CWebSocket(`ws://${this.$getHost()}/ws/fileio`)
    ws.onmessage = (e) => {
      if (typeof e.data === 'string') {
        let data = JSON.parse(event.data);
        let dtil = []
        if (data.length > 0){
          dtil.push( data[0].datetime);
        }
        for(const f of data){
          dtil.push(f.read_bytes_per_sec);
        }
        this.dtil.push(dtil);
        if (this.dtil.length > 10) {
          this.dtil.shift();
        }
      }
    }
  }, 
  data() {
    return {
      header: [],
      dtil: [],
      chartOptions: {
        title: 'test',
        subtitle: 'read'
      }
    }
  }
}
</script>

<style lang="scss">
.file-io-table{
  display:flex;
  p{
    margin-bottom: 2px;
  }
}
.file-name{
  width:75%;
}
.file-io{
  width:25%;
}
</style>

