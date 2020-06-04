<template>
  <v-col cols=12 :md="size" >
    <v-card 
      class="px-4 py-2" 
      dark>
      <v-card-title >
       Read
      </v-card-title>
        <line-chart :chart-data="chartData" :options="chartOptions" :height="75 * size" /> 
    </v-card>
  </v-col>
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
        for(const f of res){

          this.datasets.push(
            {
              label: f.file_name + "_input",
              fill: false,
              pointRadius: 0,
              borderWidth:2,
              lineTension:0,
              data: []
            }
          )
          this.datasets.push(
            {
              label: f.file_name + "_output",
              fill: false,
              pointRadius: 0,
              borderWidth:2,
              lineTension:0,
              data: []
            }
          )
        }
      })

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
      refresh: false,
    }
  }
}
</script>
