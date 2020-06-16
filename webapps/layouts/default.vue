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
            <v-tooltip bottom>
              <template v-slot:activator="{on,attrs}">
          <v-list-item-title v-bind="attrs" v-on="on">
              {{instanceName}}
          </v-list-item-title>
              </template>
              <span>{{instanceName}}</span>
            </v-tooltip>
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
      ...mapGetters(["instanceName"])
    },
}
</script>

<style scoped>
.logo {
  width: 28px;
  height: 28px;
}
</style>
