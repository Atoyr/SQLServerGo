<template>
  <v-app dark :style="{background: $vuetify.theme.themes['dark'].background}">
    <h1 class="left-label">{{ServerName}}</h1>
      <nuxt />
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
      ...mapGetters('database',["ServerName"])
    },
}
</script>

<style>

.logo {
  width: 28px;
  height: 28px;
}

</style>
