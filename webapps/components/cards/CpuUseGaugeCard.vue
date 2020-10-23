<template>
  <v-card 
    color="primary"
    dark>
    <v-card-title class="title accent py-1">
     CPU
    </v-card-title>
    <v-card-text class="py-2">
      <Gauge :value="CpuUseValue" size="100" stroke-width="8" caution-value="80" warning-value="90" unit="%" />
    </v-card-text>
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
