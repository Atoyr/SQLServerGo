<template>
  <v-card 
    color="primary"
    elevation="8"
    tile
    dark>
    <v-card-title class="subtitle-1 accent py-0 ">
      {{this.database}} TableSize
    </v-card-title>
    <v-card-text class="py-2 d-flex align-center" >
      <GChart
        type="BarChart"
        :data="chartData"
        :options="chartOptions"
        />
    </v-card-text>
  </v-card>
</template>


<script>
import { GChart } from 'vue-google-charts';
export default {
  components: {
    GChart
  },
  computed: {
    chartData: function() {
      return [this.header].concat(this.dtil);
    },
    chartOptions:function(){
      return {
        isStacked: true,
        height: this.dtil.length * 12 + 32,
        colors: ["#FFBB78","#AEC7E8","#98E28A","#FF9896","#C5B0D5","#C49C94"],
        backgroundColor:{
          fill: '#293349'
        },
        bar: { groupWidth: '50%' },
        bars: 'horizontal',
        vAxis: {
          textStyle: {
            fontSize: 8,
            color: "#eeeeee"
          }
        },
        chartArea: {
          backgroundColor: '#293349',
          left:200,
          height: this.dtil.length * 12
        },
      }
    },
  },
  filters:{
    toLocaleString(value){ 
      return value.toLocaleString();
    },
  },
  props: {
    database: {
      default: 'master'
    },
  },
  mounted: function mounted() {
    this.getTabelSize()
  },
  created() {
    this.unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === 'database/updateTableSize') {
        this.dtil = []
        for (const ts of state.database.tableSizes[this.database]){
          this.dtil.push([ts.table_name, ts.data_bytes, ts.index_bytes, ts.unused_bytes])
        }
      }
    })
  },
  beforeDestroy() {
    this.unsubscribe();
  },
  methods: {
    getTabelSize() {
       this.$store.dispatch('database/fetchTableSize',{database: this.database});
    }
  },
  data() {
    return {
      unsubscribe: {},
      header: ["tableName", "dataBytes", "indexBytes", "unusedBytes"],
      dtil : [],
    }
  }
}
</script>



