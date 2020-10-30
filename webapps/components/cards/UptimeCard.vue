<template>
  <v-card 
    color="primary"
    elevation="8"
    tile
    dark>
    <v-card-title class="subtitle-1 accent py-0 ">
     Uptime
    </v-card-title>
    <v-card-text class="py-2 d-flex align-center" style="min-height:123px">
      <v-layout justify-center>
        <h2>
        {{day}} : {{hour | zeroPadding}} : {{minute | zeroPadding}} : {{second | zeroPadding}}
        </h2>
      </v-layout>
    </v-card-text>
  </v-card>
</template>


<script>
import { mapGetters} from 'vuex'
export default {
  filters:{
    zeroPadding(value){ 
      return value.toString().padStart(2, 0);
    },
  },
  computed: {
    day: {
      get:function() {
        return Math.floor(this.diffTime / 1000 / 60 / 60 / 24);
      }
    },
    hour: {
      get:function() {
        return Math.floor(this.diffTime / 1000 / 60 / 60) % 24;
      }
    },
    minute: {
      get:function() {
        return Math.floor(this.diffTime / 1000 / 60) % 60;
      }
    },
    second: {
      get:function() {
        return Math.floor(this.diffTime / 1000) % 60;
      }
    }

              // second = 
     // Uptimestr: {
     //   get: function () {
     //     console.log(this.$store)
     //   //  let sTime = this.$store.database.state.StartTime;
     //   //  console.log("stime is",stime)
     //     return ""
     //   }
     // }
  },
  props: {
  },
  mounted: function mounted() {
    console.log(this)
    this.startTime = this.$store.getters['database/StartTime'];
    this.progress()
  },
  methods: {
    progress : function(){
      let now = Date.now();
      this.diffTime = now - this.startTime
      this.animationId = window.requestAnimationFrame(this.progress);
    }
  },
  data() {
    return {
      startTime: {},
      diffTime: 0,
      animationId: {}
    }
  }
}
</script>

