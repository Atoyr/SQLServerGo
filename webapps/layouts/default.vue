<template>
  <v-app dark :style="{background: $vuetify.theme.themes['dark'].background}">
    <v-app-bar class="accent" flat>
      <v-app-bar-title>{{ServerName}}</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-select
        label="filter"
        v-model="databaseFilter"
        :items="Databases"
        multiple
        chips
        class="align-center"
        ></v-select>
    </v-app-bar>
    <v-content>
      <nuxt />
    </v-content>
  </v-app>
</template>

<script>
import { mapGetters} from 'vuex'
export default {
  data() {
    return {
      dbname : "",
      drawer: true,
      mini: true,
    }
  },
  computed: {
    ...mapGetters('database',["ServerName","Databases"]),
    databaseFilter: {
      get () { return this.$store.getters["database/DatabaseFilter"]; },
      set (val) { this.$store.commit('database/updateDatabaseFilter', {databaseFilter: val}) },
    },
  },
}
</script>

<style>

.logo {
  width: 28px;
  height: 28px;
}

</style>
