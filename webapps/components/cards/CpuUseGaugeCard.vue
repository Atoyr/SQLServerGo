<template>
  <v-card 
    class="py-2" 
    dark>
    <h1 class="label">
     CPU (%)
    </h1>
    <div class="flex">
      <Gauge :value="CpuUseValue" size="300" stroke-width="16" caution-value="80" warning-value="90" unit="%" />
      <div class="flex-col px-2">
        <h2 class="cpu-use-label">sqlserer use</h2>
        <h3 class="py-2">{{this.sqlUseValue}} %</h3>
        <h2 class="cpu-use-label">other use</h2>
        <h3 class="py-2">{{this.otherUseValue}} %</h3>
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
  computed: {
    CpuUseValue: {
      get: function () {
        return this.sqlUseValue + this.otherUseValue;
      }
    }
  },
  methods: {
  },
  props: {
  },
  mounted() {
    fetch(`http://${this.$getHost()}/api/cpuUsed`)
      .then(response => {
        return response.json();
      })
      .then(res => {
        this.sqlUseValue = res.sql_process_utilization;
        this.otherUseValue = res.other_process_utilization;
      })

     var cpuws = new W3CWebSocket(`ws://${this.$getHost()}/ws/cpu`)
     cpuws.onmessage = (e) => {
        let data = JSON.parse(event.data);
        this.sqlUseValue = data.sql_process_utilization;
        this.otherUseValue = data.other_process_utilization;
     }
  }, 
  data() {
    return {
      sqlUseValue:0,
      otherUseValue:0,
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
