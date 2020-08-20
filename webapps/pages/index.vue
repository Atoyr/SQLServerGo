<template>
  <div class="px-2 py-2">
    <v-row>
      <v-col md=4>
        <CpuUseGaugeCard class="px-2"></CpuUseGaugeCard>
      </v-col>
      <v-col md=4>
        <MemoryUseGaugeCard class="px-2"></MemoryUseGaugeCard>
      </v-col>
      <v-col md=4>
        <BufferCacheRateGaugeCard class="px-2"></BufferCacheRateGaugeCard>
      </v-col>
    </v-row>
    <v-row>
      <v-col md=6>
        <FileReadIOCard></FileReadIOCard>
      </v-col>
      <v-col md=6>
        <FileWriteIOCard></FileWriteIOCard>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import CpuUseGaugeCard from '@/components/cards/CpuUseGaugeCard.vue'
import MemoryUseGaugeCard from '@/components/cards/MemoryUseGaugeCard.vue'
import BufferCacheRateGaugeCard from '@/components/cards/BufferCacheRateGaugeCard.vue'
import FileReadIOCard from '@/components/cards/FileReadIOCard.vue'
import FileWriteIOCard from '@/components/cards/FileWriteIOCard.vue'
import Gauge from '@/components/Gauge.vue'
import axios from 'axios'

export default {
  components: {
    Gauge,
    MemoryUseGaugeCard,
    CpuUseGaugeCard,
    BufferCacheRateGaugeCard,
    FileReadIOCard,
    FileWriteIOCard
  },
   async fetch ({ store, params }) {
     console.log("fetch")
     await store.dispatch('database/fetchServerProperty');
     await store.dispatch('database/fetchDatabaseFiles');
   },
}
</script>
<style lang="scss">
.test {
  text-align: center;

}
</style>
