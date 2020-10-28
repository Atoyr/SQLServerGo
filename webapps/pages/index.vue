<template>
  <div class="px-2 py-2">
    <v-row>
      <v-col md=3>
        <Uptime ></Uptime>
      </v-col>
      <v-col md=2>
        <CpuUseGaugeCard ></CpuUseGaugeCard>
      </v-col>
      <v-col md=2>
        <MemoryUseGaugeCard ></MemoryUseGaugeCard>
      </v-col>
      <v-col md=2>
        <BufferCacheRateGaugeCard ></BufferCacheRateGaugeCard>
      </v-col>
    </v-row>
    <v-row>
      <!--
      <v-col md=6>
        <FileReadIOCard></FileReadIOCard>
      </v-col>
      <v-col md=6>
        <FileWriteIOCard></FileWriteIOCard>
      </v-col>
      -->
    </v-row>
  </div>
</template>

<script>
import CpuUseGaugeCard from '@/components/cards/CpuUseGaugeCard.vue'
import MemoryUseGaugeCard from '@/components/cards/MemoryUseGaugeCard.vue'
import BufferCacheRateGaugeCard from '@/components/cards/BufferCacheRateGaugeCard.vue'
import Uptime from '@/components/cards/Uptime.vue'
import FileReadIOCard from '@/components/cards/FileReadIOCard.vue'
import FileWriteIOCard from '@/components/cards/FileWriteIOCard.vue'
import Gauge from '@/components/Gauge.vue'
import axios from 'axios'
import { mapGetters} from 'vuex'

export default {
  components: {
    MemoryUseGaugeCard,
    CpuUseGaugeCard,
    BufferCacheRateGaugeCard,
    Uptime,
    FileReadIOCard,
    FileWriteIOCard
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
