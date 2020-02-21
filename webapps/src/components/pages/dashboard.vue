<template>
  <div>
    <line-chart :chartdata="chartdata" :options="option" ref="testchart"></line-chart>
    <button v-on:click="update()"></button>
  </div>
</template>
<script>
import { Line } from 'vue-chartjs'
import LineChart from '@/components/modules/LineChart.vue';
import 'chartjs-plugin-streaming';
import Vue from 'vue';
import { w3cwebsocket } from 'websocket'
const W3CWebSocket = w3cwebsocket

export default {
  name: 'dashboard',
  components: {
    LineChart
  },
  data ()  {
    return {
      chartdata: {},
      option: null,
      socket: new W3CWebSocket('ws://localhost:8080/ws')
    }
  },
  mounted () {
    this.chartdata = {
      datasets: [{
        label: 'Dataset 1',
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.5)',
        lineTension: 0,
        borderDash: [8, 4],
        data: []
      }]
    }
    this.option = {
      scales: {
        xAxes: [{
          type: 'realtime',
          realtime: {
            delay: 2000,
            refresh: 1000
          }
        }],
        yAxes: [{
          scaleLabel: {
            display: true,
            labelString: 'value'
          }
        }]
      },
      preservation: false
    }

    this.socket.onmessage = function(event) {
      let data = JSON.parse(event.data)
      console.log(data)
    }
    this.$refs.testchart.addDataset(this.chartdata.datasets[0])
    this.$refs.testchart.applyOption(this.option)
  },
  methods : {
    update: function() {
      this.$refs.testchart.onReceive({
        index: 0,
        timestamp: Date.now(),
        value: 100
      })
    },
  }
}
</script>
