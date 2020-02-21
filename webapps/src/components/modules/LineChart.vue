<script>
import { Line } from 'vue-chartjs'
import 'chartjs-plugin-streaming';

export default {
  extends: Line,
  props: {
    chartdata: {
      type: Object,
      default: null
    },
    options: {
      type: Object,
      default: null
    }
  },
  data ()  {
    return {
      cd: {},
      op: null,
    }
  },
  mounted () {
    console.log(this)
    this.cd = {datasets: [ ]}
    //this.renderChart(this.chartdata, this.options)
    this.renderChart(this.cd,this.op)
    this.$data._chart.update({
      preservation: true
    });
  },
  methods: {
    update: function () {
      console.log(this.chartdata)
      this.$data._chart.update({
        preservation: true
      });
    },
    addDataset: function(dataset) {
      this.cd.datasets.push(dataset)
      this.renderChart(this.cd,this.op)
    },
    applyOption: function(option) {
      this.op = option
      this.renderChart(this.cd,this.op)
    },

    onReceive: function (event) {
      console.log(event)
      this.cd.datasets[event.index].data.push({
        x: event.timestamp,
        y: event.value
      });
      this._data._chart.update({
        preservation: true
      });
    }
  },
  plugins: {
    streaming: {
        frameRate: 30
    }
  }
} 
</script>
