<template>
  <v-card 
    color="primary"
    tile
    dark>
    <v-card-title class="subtitle-1 accent py-0">
      {{this.database}} File {{this.write?'Write':'Read'}} KiB/sec
    </v-card-title>
    <v-card-text>
      <v-container class="px-1 py-1">
        <v-row>
          <v-col md=6>
            <GChart
              type="AreaChart"
              :data="chartData"
              :options="chartOptions"
              />
          </v-col>         
          <v-col md=6>
            <v-simple-table class="primary" dense dark height=150>
                <thead>
                  <tr>
                    <th class="text-left">
                      fileName
                    </th>
                    <th class="text-left">
                      min
                    </th>
                    <th class="text-left">
                      avg
                    </th>
                    <th class="text-left">
                      max
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="item in tableData"
                    :key="item.name"
                  >
                  <td class="pa-1 text-left body-2">
                    <v-icon :color="item.color">mdi-chart-line-variant</v-icon>
                    {{ item.name }}
                  </td>
                  <td class="pa-1 text-right body-2">{{ item.min | round }}KiB/s</td>
                  <td class="pa-1 text-right body-2">{{ item.avg | round }}KiB/s</td>
                  <td class="pa-1 text-right body-2">{{ item.max | round }}KiB/s</td>
                  </tr>
                </tbody>
            </v-simple-table>
          </v-col>
        </v-row>
      </v-container>
    </v-card-text>
  </v-card>
</template>

<script>
import { GChart } from 'vue-google-charts';
import { mapGetters } from 'vuex';

export default {
  components: {
    GChart
  },
  filters:{
    round(value){ 
      return Math.round(value * 100) / 100;
    },
  },
  methods: {
  },
  computed: {
    ...mapGetters('database',["Instance"]),
    chartData: function() {
      return this.header.concat(this.dtil);
    },
  },
  created() {
    this.unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === 'database/updateInstance') {
        let instance = state.database.instance
        if (this.database in instance){
          let h = instance[this.database].files.concat();
          h.unshift("datetime");
          this.header = [h];
          let d = this.write ? instance[this.database].write: instance[this.database].read;
          this.dtil = d.concat();
        }

        let tdata = []
        let index = this.dtil.length - 1
        let summary = state.database.instance[this.database].summary
        for(let i = 0; i < summary.length; i++) {

          tdata.push({
            name: summary[i].fileName,
            min: this.write ? summary[i].MinWrite : summary[i].MinRead,
            avg: this.write ? summary[i].AvgWrite : summary[i].AvgRead,
            max: this.write ? summary[i].MaxWrite : summary[i].MaxRead,
            color:this.chartOptions.colors[i % this.chartOptions.colors.length] })
        }
        this.tableData = tdata
      }
    });
  },
  beforeDestroy() {
    this.unsubscribe();
  },
  props: {
    size: Number,
    database: {
      default: 'master'
    },
    write: Boolean
  },
  mounted() {
  }, 
  data() {
    return {
      header: [],
      dtil : [],
      tableData : [],
      chartOptions: {
        height: 150,
        colors: ["#FFBB78","#AEC7E8","#98E28A","#FF9896","#C5B0D5","#C49C94"],
        hAxis: {
          format: 'hh:mm:ss',
          textStyle: {
            color: '#eeeeee'
          }
        },
        vAxis: {
          minValue: 100,
          textStyle: {
            color: '#eeeeee'
          }
        },
        backgroundColor:{
          fill: '#293349'
        },
        chartArea: {
          backgroundColor: '#293349',
          width: "90%",
          top:8,
          left:32,
          bottom: 32,
        },
        legend: {
          position: "none"
        }
      },
    }
  }
}
</script>

<style lang="scss">
.file-io-table{
  display:flex;
  p{
    margin-bottom: 2px;
  }
}
.file-name{
  width:75%;
}
.file-io{
  width:25%;
}
</style>


