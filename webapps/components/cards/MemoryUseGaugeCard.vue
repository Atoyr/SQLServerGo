<template>
  <v-card 
    class="py-2" 
    dark>
    <h1 class="label">
     Memory (%)
    </h1>
    <div class="flex">
      <Gauge :value="MemoryUseValue" size="300" stroke-width="16" caution-value="80" warning-value="90" unit="%" />
      <div class="flex-col px-2">
        <h2 class="memory-use-label">Physical Memory</h2>
        <h3 class="py-2">{{this.physicalMemory}} MB</h3>
        <h2 class="memory-use-label">Used Memory</h2>
        <h3 class="py-2">{{this.usedMemory}} MB</h3>
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
    MemoryUseValue: {
      get: function () {
        return parseInt(100 * this.usedMemory / this.physicalMemory,10)
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
        this.physicalMemory = data.physical_memory;
        this.usedMemory = data.used_memory;
        this.avaliableMemory = data.avaliable_memory;
        this.memoryStatus = data.system_memory_state;
      })

     var memoryws = new W3CWebSocket(`ws://${this.$getHost()}/ws/memory`)
     memoryws.onmessage = (e) => {
        let data = JSON.parse(event.data);
        this.physicalMemory = data.physical_memory;
        this.usedMemory = data.used_memory;
        this.avaliableMemory = data.avaliable_memory;
        this.memoryStatus = data.system_memory_state;
     }
  }, 
  data() {
    return {
      physicalMemory:0,
      usedMemory:0,
      avaliableMemory:0,
      memoryStatus:""
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
.memory-use-label{
  border-left-style: solid;
  padding: 4px;
}
</style>
