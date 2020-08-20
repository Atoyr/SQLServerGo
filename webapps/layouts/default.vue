<template>
  <v-app dark>
    <h1 class="left-label">{{InstanceName}}</h1>
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
      ...mapGetters('database',["InstanceName"])
    },
}
</script>

<style>

.logo {
  width: 28px;
  height: 28px;
}

</style>
