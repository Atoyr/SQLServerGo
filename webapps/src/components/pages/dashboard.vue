<template>
  <div>
    <line-chart class="chart" ref="readIOChart" :height="540"></line-chart>
    <line-chart class="chart" ref="writeIOChart" :height="540"></line-chart>
  </div>
</template>
<script>
import { Line } from 'vue-chartjs';
import LineChart from '@/components/modules/LineChart.vue';
import 'chartjs-plugin-streaming';
import Vue from 'vue';
import { w3cwebsocket } from 'websocket';
import 'chartjs-plugin-colorschemes';
import { RedGold21 } from 'chartjs-plugin-colorschemes/src/colorschemes/colorschemes.tableau'

const W3CWebSocket = w3cwebsocket;

export default {
  name: 'dashboard',
  components: {
    LineChart
  },
  data ()  {
    return {
      socket: new W3CWebSocket(`ws://${window.location.host}/ws`),
      isPause: false
    }
  },
  mounted () {
    fetch(`${window.location.origin}/api/databaseFiles`)
    .then(response => {
      return response.json();
    })
    .then(res => {
      for(const f of res){
        this.$refs.readIOChart.addDataset(
          {
            label: f.file_name,
            data: []
          }
        );
        this.$refs.writeIOChart.addDataset(
          {
            label: f.file_name,
            data: []
          }
        );
      }
    })
    let option = {
      scales: {
        xAxes: [{
          type: 'realtime',
          realtime: {
            duration: 40000,
            delay: 5000,
            ttl: 60000,
            pause: this.isPause
          },
        }],
        yAxes: [{
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
      plugins: {
        colorschemes: {
          scheme: RedGold21
        }
      },
      legend: {
        position: 'right'
      },
      preservation: true,
      responsive: true,
      maintainAspectRatio: false
    }
    this.$refs.readIOChart.applyOption(option)
    this.$refs.writeIOChart.applyOption(option)

    this.socket.onmessage = (event) => {
      let data = JSON.parse(event.data)
      for(const f of data){
        this.$refs.readIOChart.onReceive({
          index: f.id - 1,
          timestamp: Date.parse(f.datetime),
          value: f.read_bytes_per_sec
        });
        this.$refs.writeIOChart.onReceive({
          index: f.id - 1,
          timestamp: Date.parse(f.datetime),
          value: f.write_bytes_per_sec
        });
      }
      this.$refs.readIOChart.update()
      this.$refs.writeIOChart.update()
    }
  },
  methods : {
  }
}
</script>
