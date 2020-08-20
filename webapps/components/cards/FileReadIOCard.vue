<template>
  <v-card 
    class="py-2" 
    dark>
    <h1 class="label">
     File Read bytes/sec
    </h1>
    <div class="px-4">
      <line-chart :chart-data="chartData" :options="chartOptions" height=300 /> 
    </div>
    <div class="file-io-table px-4">
      <div class="file-name">
        <p v-for="fileName in files">
        {{fileName}}
        </p>
      </div>
      <div class="file-io">
        <p v-for="readValue in readValues">
        {{readValue}} bytes/sec
        </p>
      </div>
    </div>

  </v-card>
</template>

<script>
import { mapState } from 'vuex';
import colors from 'vuetify/es5/util/colors';
import { w3cwebsocket } from 'websocket';
import 'chartjs-plugin-streaming';
import 'chartjs-plugin-colorschemes';
import { SetOne5 } from 'chartjs-plugin-colorschemes/src/colorschemes/colorschemes.brewer'
const W3CWebSocket = w3cwebsocket
const datalength = 300

export default {
  methods: {
  },
  computed: {
    chartData: function() {
      let labels = this.labels
      let datasets = this.datasets
      let refresh = this.refresh
      this.refresh = !refresh
      return {
        labels: labels,
        datasets: datasets,
      }
    },
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
        for(const f of res){

          this.datasets.push(
            {
              label: f.file_name,
              fill: false,
              pointRadius: 0,
              borderWidth:2,
              lineTension:0,
              data: []
            }
          )
          this.files.push(f.file_name);
        }
      })

    var ws = new W3CWebSocket(`ws://${this.$getHost()}/ws/fileio`)
    ws.onmessage = (e) => {
      if (typeof e.data === 'string') {
        let data = JSON.parse(event.data);
        if (data.length > 0){
          this.labels.push( Date.parse(data[0].datetime));
        }
        this.readValues.length=0;
        for(const f of data){
          let id = f.id -1

          this.datasets[id].data.push({
            x: Date.parse(f.datetime),
            y: f.read_bytes_per_sec,
          });

          if (this.datasets[id].data.length > 120) {
            this.datasets[id].data.pop();
          }
          this.readValues.push(f.read_bytes_per_sec);
        }
        if (this.labels.length > 120) {
          this.labels.pop();
        }
      }
    }

//     var cpuws = new W3CWebSocket(`ws://${this.$getHost()}/ws/cpu`)
//     cpuws.onmessage = (e) => {
//       console.log(e.data)
//     }
 
    this.refresh = true

  }, 
  data() {
    return {
      chartOptions: {
        scales: {
          xAxes: [{
            type: 'realtime',
            gridLines: {
              color: "#6d6d6d"
            },
            scaleLabel: {
              display: false,
            },
            realtime: {
              duration: 40000,
              delay: 2000,
              ttl: 60000,
              pause: false,
            },
          }],
          yAxes: [{
            gridLines: {
              color: "#6d6d6d"
            },
            scaleLabel: {
              display: true,
              labelString: 'value'
            },
            ticks:{
              beginAtZero: true,
              min:0
            }
          }]
        },
        legend: {
          display: false,
          position: 'right'
        },
        preservation: true,
        responsive: true,
        maintainAspectRatio: false,
        padding: {
          left: 50,
          right: 50,
          top: 50,
          bottom: 51
        },
        plugins: {
          colorschemes: {
            scheme: SetOne5
          }
        },
      },
      labels: [],
      datasets: [],
      files: [],
      readValues: [],
      refresh: false,
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
