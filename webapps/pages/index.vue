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
         await store.dispatch('fetchInstance');
       },
}
</script>
