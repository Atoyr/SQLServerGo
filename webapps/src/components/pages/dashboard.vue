<template>
  <line-chart :chart-data="readDataset"></line-chart>
  <line-chart :chart-data="writeDataset"></line-chart>
</template>

<script>
import LineChart from '@/components/modules/LineChart.vue';
import 'chartjs-plugin-streaming';

export default {
  name: 'dashboard',
  components: {
    LineChart
  },
  data: () => ({

  }),
  mounted () {

  }
}
</script>
