import Vue from 'vue';
import { Line, mixins } from 'vue-chartjs';
const { reactiveProp } = mixins;

Vue.component('line-chart', {
  extends: Line,
  mixins: [reactiveProp],
  props: {
    chartdata: {
      type: Object,
      default: null,
    },
    options: {
      type: Object,
      default: null,
    },
  },
  mounted() {
    this.renderChart(this.chartData, this.options);
  },
});

