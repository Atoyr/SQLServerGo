<template>
  <div>
    <line-chart :chartdata="chartdata" :options="option" ref="testchart"></line-chart>
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
      socket: new W3CWebSocket(`ws://${window.location}/ws`)
    }
  },
  mounted () {
//    this.chartdata = {
//      datasets: []
//    };

    fetch(`http://${window.location}/api/databaseFiles`)
      .then(response => {
        return response.json();
      })
      .then(res => {
        for(const f of res){
          this.$refs.testchart.addDataset(
            {
              label: f.file_name,
              data: []
            }
          );
        }
      })
//    this.chartdata = {
//      datasets: [{
//        label: 'Dataset 1',
//        borderColor: 'rgb(255, 99, 132)',
//        backgroundColor: 'rgba(255, 99, 132, 0.5)',
//        lineTension: 0,
//        borderDash: [8, 4],
//        data: []
//      }]
//    }
    this.$refs.testchart.applyOption({
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
      preservation: true
    });

    this.socket.onmessage = function(event) {
      let data = JSON.parse(event.data)
      for(const f of data){
        this.$refs.testchart.onReceive({
          index: f.id,
          timestamp: f.datetime,
          value: f.read_bytes_per_sec
        })
      }
    }
  },
  methods : {
  //  update: function() {
  //    this.$refs.testchart.onReceive({
  //      index: 0,
  //      timestamp: Date.now(),
  //      value: 100
  //    })
  //  },
  }
}
</script>
