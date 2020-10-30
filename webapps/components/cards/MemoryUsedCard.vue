<template>
  <v-card 
    color="primary"
    elevation="8"
    tile
    dark>
    <v-card-title class="subtitle-1 accent py-0">
     Memory
    </v-card-title>
    <v-card-text class="py-2">
      <GChart
        type="AreaChart"
        :data="chartData"
        :options="chartOptions"
        />
    </v-card-text>
  </v-card>
</template>


<script>
import { w3cwebsocket } from 'websocket';
const W3CWebSocket = w3cwebsocket
import { GChart } from 'vue-google-charts';

export default {
  components: {
    GChart,
  },
  computed: {
    chartData: function() {
      return this.header.concat(this.dtil);
    },
    chartOptions: function(){
      return {
      // title: this.write ? 'write' : 'read',
        colors: ["#FFBB78","#AEC7E8"],
        hAxis: {
          format: 'hh:mm:ss',
          textStyle: {
            color: '#eeeeee'
          }
        },
        vAxis: {
          minValue: this.physicalMemory,
          textStyle: {
            color: '#eeeeee'
          }
        },
        backgroundColor:{
          fill: '#293349'
        },
        isStacked: "true",
        chartArea: {
          backgroundColor: '#293349',
        },
        legend: {
          position: "none"
        }
      }
    },
  },
  methods: {
  },
  props: {
  },
  mounted() {
    fetch(`http://${this.$getHost()}/api/memory`)
      .then(response => {
        return response.json();
      })
      .then(res => {
        this.physicalMemory = res.physical_memory;
      })

     var memoryws = new W3CWebSocket(`ws://${this.$getHost()}/ws/memory`)
     memoryws.onmessage = (e) => {
        let data = JSON.parse(event.data);
        this.physicalMemory = data.physical_memory;
        this.dtil.push([new Date(data.datetime), data.used_memory, data.available_memory])
        if (this.dtil.length > 30) {
          this.dtil.shift();
        }
     }
  }, 
  data() {
    return {
      header: [["datetime","useMemory","availableMemory"]],
      dtil : [],
      physicalMemory : 8192
    }
  }
}
</script>

