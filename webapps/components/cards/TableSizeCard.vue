<template>
  <v-card 
    color="primary"
    elevation="8"
    tile
    dark>
    <v-card-title class="subtitle-1 accent py-0 ">
      {{this.database}} TableSize
    </v-card-title>
    <v-card-text class="py-2 d-flex align-center" style="min-height:123px">
      <v-data-table
        :headers="headers"
        :items="tableSizes"
        :items-per-page="10"
        dense
        class="primary"
        ></v-data-table>
    </v-card-text>
  </v-card>
</template>


<script>
export default {
  computed: {
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
        this.tableSizes = state.database.tableSizes[this.database]
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
      headers: [
        { text: 'tableName', value: 'table_name', width: 300},
        { text: 'rows', value: 'rows' },
        { text: 'reservedBytes', value: 'reserved_bytes' },
        { text: 'dataBytes', value: 'data_bytes' },
        { text: 'indexBytes', value: 'index_bytes' },
        { text: 'unusedBytes', value: 'unused_bytes' },
      ],
      unsubscribe: {},
      tableSizes:[]
    }
  }
}
</script>


