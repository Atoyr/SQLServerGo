<template>
  <v-card 
    class="py-2" 
    dark>
    <h1 class="label">
     Buffer Cache Hit Rate (%)
    </h1>
    <div class="flex">
      <Gauge :value="Rate" size="300" stroke-width="16" caution-value="80" warning-value="90" unit="%" 
      :stroke-colors="strokeColors" :stroke-bg-colors="strokeBgColors"/>
      <div class="flex-col px-2">
        <h2 class="cpu-use-label">Buffer Cache Hit Rate</h2>
        <h3 class="py-2">{{this.Rate}} %</h3>
      </div>
    </div>
  </v-card>
</template>


<script>
import { w3cwebsocket } from 'websocket';
const W3CWebSocket = w3cwebsocket
import Gauge from '@/components/Gauge.vue'

export default {
  components: {
    Gauge,
  },
  methods: {
  },
  computed: {
    Rate: {
      get: function () {
        return parseInt(this.rate)
      }
    }
  },
  props: {
  },
  mounted() {
    fetch(`http://${this.$getHost()}/api/buferCache`)
      .then(response => {
        return response.json();
      })
      .then(res => {
        this.rate = res.buffer_cache_rate
      })

     var bufcachews = new W3CWebSocket(`ws://${this.$getHost()}/ws/bufferCache`)
     bufcachews.onmessage = (e) => {
        let data = JSON.parse(event.data);
        this.rate = data.buffer_cache_rate
     }
  }, 
  data() {
    return {
      rate:0,
      strokeColors: ['#ff7f7f','#ffb266','#7fff7f'],
      strokeBgColors: ['#ffcccc','#ffd8b2','#ccffcc']
    }
  }
}
</script>

<style lang="scss">
.flex{
  display: flex;
  flex-direction: row;
}
.flex-col{
  display: inline-flex;
  flex-direction: column;
  text-align: left;
}
.cpu-use-label{
  border-left-style: solid;
  padding: 4px;
}
</style>

