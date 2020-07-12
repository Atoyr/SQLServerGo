<template>
  <div>
    <InputChart size=4 />
  </div>
</template>

<script>
import InputChart from '@/components/cards/inputChart.vue'
import axios from 'axios'

export default {
  components: {
    InputChart
  },
   async fetch ({ store, params }) {
     console.log("fetch")
     await store.dispatch('database/fetchServerProperty');
   },
}
</script>
