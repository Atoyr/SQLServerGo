<template>
  <v-card 
    color="primary"
    tile
    dark>
    <v-card-title class="title accent py-1">
      {{this.database}} File {{this.write?'write':'read'}} MiB/sec
    </v-card-title>
    <v-card-text>
      <v-container class="px-1 py-1">
        <v-row>
          <v-col md=8>
            <GChart
              type="LineChart"
              :data="chartData"
              :options="chartOptions"
              />
          </v-col>         
          <v-col md=4>
            <v-simple-table class="primary" dense dark>
                <thead>
                  <tr>
                    <th class="text-left">
                      fileName
                    </th>
                    <th class="text-left">
                      KiB/sec
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="item in tableData"
                    :key="item.name"
                  >
                  <td class="pa-1">{{ item.name }}</td>
                  <td class="pa-1">{{ item.io }} KiB/sec</td>
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
        for(let i = 1; i < this.header[0].length; i++) {
          tdata.push({name: this.header[0][i],io: this.dtil[index][i]})
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
        title: this.write ? 'write' : 'read',
        colors: ["#AEC7E8","#FFBB78","#98E28A","#FF9896","#C5B0D5","#C49C94"],
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


