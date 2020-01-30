<script>
import { Line, mixins } from 'vue-chartjs'
import 'chartjs-plugin-streaming';
const { reactiveProp} = mixins

export default {
  extends: Line,
  mixins: [reactiveProp],
  props: ['options'],
  method() {
    this.renderChart(this.chartData, this.options)
  }
}
</script>
