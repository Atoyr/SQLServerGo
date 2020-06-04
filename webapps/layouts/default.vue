<template>
  <v-app dark>
    <v-card
      class="toolbar-card"
      color="grey lighten-4"
      flat
      tile
      >
      <v-toolbar dense flat class="toolbar">
        <v-app-bar-nav-icon>
        </v-app-bar-nav-icon>
        <v-toolbar-title>{{dbname}}</v-toolbar-title>
      <v-spacer></v-spacer>
      
      </v-toolbar>
    </v-card>
    <v-container fluid>
      <nuxt />
    </v-container>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      dbname : "",
    }
  }
}
</script>
