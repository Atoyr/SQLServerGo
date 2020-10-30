<template>
  <v-container fluid class="px-4 py-1">
    <v-row>
      <v-col sm=12 md=3 class="pa-1">
        <UptimeCard ></UptimeCard>
      </v-col>
      <v-col sm=4 md=2 class="pa-1">
        <CpuUseGaugeCard ></CpuUseGaugeCard>
      </v-col>
      <v-col sm=4 md=2 class="pa-1">
        <MemoryUseGaugeCard ></MemoryUseGaugeCard>
      </v-col>
      <v-col sm=4 md=2 class="pa-1">
        <BufferCacheRateGaugeCard ></BufferCacheRateGaugeCard>
      </v-col>
    </v-row>
    <v-row>
      <v-col md=6 class="pa-1">
        <FileIOCard database="R_1_1_0_SC" ></FileIOCard>
      </v-col>
      <v-col md=6 class="pa-1">
        <FileIOCard database="R_1_1_0_SC" write></FileIOCard>
      </v-col>
    </v-row>
    <v-row>
      <v-col md=6 class="pa-1">
        <FileIOCard database="R_1_1_0_FI" ></FileIOCard>
      </v-col>
      <v-col md=6 class="pa-1">
        <FileIOCard database="R_1_1_0_FI" write></FileIOCard>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import CpuUseGaugeCard from '@/components/cards/CpuUseGaugeCard.vue'
import MemoryUseGaugeCard from '@/components/cards/MemoryUseGaugeCard.vue'
import BufferCacheRateGaugeCard from '@/components/cards/BufferCacheRateGaugeCard.vue'
import UptimeCard from '@/components/cards/UptimeCard.vue'
import FileReadIOCard from '@/components/cards/FileReadIOCard.vue'
import FileWriteIOCard from '@/components/cards/FileWriteIOCard.vue'
import testChart from '@/components/cards/testChart.vue'
import FileIOCard from '@/components/cards/FileIOCard.vue'
import Gauge from '@/components/Gauge.vue'
import { mapGetters} from 'vuex'
import { w3cwebsocket } from 'websocket';
const W3CWebSocket = w3cwebsocket
const datalength = 300

export default {
  components: {
    MemoryUseGaugeCard,
    CpuUseGaugeCard,
    BufferCacheRateGaugeCard,
    UptimeCard,
    FileReadIOCard,
    FileWriteIOCard,
    FileIOCard,
    testChart
  },
  mounted() {
    var ws = new W3CWebSocket(`ws://${this.$getHost()}/ws/fileio`)
    ws.onmessage = (e) => {
      if (typeof e.data === 'string') {
        let data = JSON.parse(event.data);
        this.$store.commit('database/updateInstance',{data})
      }
    }
  }, 
  computed: {
      ...mapGetters('database',["ServerName"])
    },
   async fetch ({ store, params }) {
     console.log("fetch")
     await store.dispatch('database/fetchServerProperty');
     await store.dispatch('database/fetchDatabaseFiles');
     await store.dispatch('database/fetchServerStatus');
   },
  head() {
    return {
      title: this.ServerName
    }
  }
}
</script>
<style lang="scss">
.test {
  text-align: center;

}
</style>
