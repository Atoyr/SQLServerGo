<template>
  <v-app dark>
    <v-card >
      <v-navigation-drawer
        v-model="drawer"
        :mini-variant.sync="mini"
        permanent
        fixed
        app
        temporary
        >
        <v-list-item class="px-2">
          <v-list-item-avatar>
            <img src="@/static/logo-sql-server.png" class="logo" @click.stop="mini = !mini"></img>
          </v-list-item-avatar>
          <v-list-item-title>{{dbname}} fooo</v-list-item-title>
          <v-btn icon @click.stop="mini = !mini">
            <v-icon>mdi-chevron-left</v-icon>
          </v-btn>
        </v-list-item>
        <v-divider></v-divider>
      </v-navigation-drawer>
    </v-card>
    <v-content>
    <v-container fluid>
      <nuxt />
    </v-container>
    </v-content>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      dbname : "",
      drawer: true,
      mini: true,
    }
  },
  fetch({ store, params}) {
    return axios.get(`http://${this.$getHost()}/api/instance`)
    .then((res) => {
      console.log(res)
    })
  }
}
</script>

<style scoped>
.logo {
  width: 28px;
  height: 28px;
}
</style>
