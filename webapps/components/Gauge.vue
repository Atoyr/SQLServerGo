<template>
  <div 
    class="guage">
    <svg 
      :width="size" 
      :height="size"
      :viewBox="`0 0 ${gaugeSize} ${gaugeSize}`" 
      xmlns="http://www.w3.org/2000/svg">
      <g 
        :transform="`rotate(150,${center},${center})`">
      <circle 
        class="guage__background-warning" 
        r="226" 
        :cy="center" 
        :cx="center" 
        :stroke-width="this.gaugeStrokeWidth - 4" 
        :stroke="strokeWarningBgColor" 
        stroke-dasharray="1420"
        stroke-dashoffset="473.3"
        stroke-linecap="round" 
        fill="none" />
      <circle 
        class="guage__background-caution" 
        r="226" 
        :cy="center" 
        :cx="center" 
        :stroke-width="this.gaugeStrokeWidth - 2" 
        :stroke="strokeCautionBgColor" 
        stroke-dasharray="1420"
        :stroke-dashoffset="warningOffset"
        stroke-linecap="round" 
        fill="none" />
      <circle 
        class="guage__background" 
        r="226" 
        :cy="center" 
        :cx="center" 
        :stroke-width="this.gaugeStrokeWidth" 
        :stroke="strokeBgColor" 
        stroke-dasharray="1420"
        :stroke-dashoffset="cautionOffset"
        stroke-linecap="round" 
        fill="none" />
      <circle 
        r="226" 
        :cy="center" 
        :cx="center" 
        :stroke-width="this.gaugeStrokeWidth" 
        :stroke="strokeColor" 
        stroke-dasharray="1420"
        :stroke-dashoffset="offset"
        stroke-linecap="round" 
        fill="none" />
     </g>
     <g>
       <text 
        :x="center" 
        :y="center + 40" 
        text-anchor="middle" 
        stroke-width="1" 
        fill="#FFFFFF"
        font-size="128" 
        font-family="Meiryo">
         {{visualValue}}
       </text>
     </g>
    </svg>
  </div>
</template>

<script>
export default {
  props:{
    value:{
      type: Number,
      default: 0
    },
    minValue: {
      type: Number,
      default: 0
    },
    maxValue: {
      type: Number,
      default: 100
    },
    cautionValue: {
      type: Number,
      default: 101
    },
    warningValue: {
      type: Number,
      default: 101
    },
    strokeWidth: {
      type: Number,
      default: 50
    },
    strokeColors: {
      type: Array,
      default: () => ['#7fff7f','#ffb266','#ff7f7f']
    },
    strokeBgColors: {
      type: Array,
      default: () => ['#ccffcc','#ffd8b2','#ffcccc']
    },
    size: {
      type: Number,
      default: 500
    },
    fontSize: {
      type: Number,
      default: 128
    },
    unit: {
      type: String,
      default: ''
    }
  },
  computed:{
    mag() {
      return 452 / (this.size - this.strokeWidth);
    },
    gaugeStrokeWidth() {
      return this.strokeWidth * this.mag;
    },
    gaugeSize() {
      return 452 + this.gaugeStrokeWidth;
    },
    center() {
      return this.gaugeSize / 2;
    },
    gaugeValue() {
      let v = parseInt(this.value,10)
      let min = parseInt(this.minValue,10)
      let max = parseInt(this.maxValue,10)
      if (v < min) {
        return min;
      }else if (max < v) {
        return max;
      }
      return v;
    },
    gaugeCautionValue() {
      let c = parseInt(this.cautionValue,10)
      let min = parseInt(this.minValue,10)
      let max = parseInt(this.maxValue,10)
      if (c < min) {
        return min;
      }else if (max < c) {
        return max;
      }
      return c;
    },
    gaugeWarningValue() {
      let w = parseInt(this.warningValue,10)
      let min = parseInt(this.minValue,10)
      let max = parseInt(this.maxValue,10)
      if (w < min) {
        return min;
      }else if (max < w) {
        return max;
      }
      return w;
    },
    percent() {
      return (this.gaugeValue - this.minValue) / Math.abs(this.maxValue - this.minValue);
    },
    visualPercent() {
      return this.isAnimating ? this.speedEaseoutPercent : this.percent;
    },
    visualValue() {
      let min = parseInt(this.minValue,10)
      let max = parseInt(this.maxValue,10)
      return min + parseInt((this.isDrumRoll ? this.visualPercent : this.percent) * (max - min));
    },
    offset() {
      return 1420 - (947 * this.visualPercent);
    },
    cautionOffset() {
      return 1420 - (947 * (this.gaugeCautionValue - this.minValue) / Math.abs(this.maxValue - this.minValue));
    },
    warningOffset() {
      return 1420 - (947 * (this.gaugeWarningValue - this.minValue) / Math.abs(this.maxValue - this.minValue) );
    },
    strokeColor() {
      let v = parseInt(this.value,10)
      let c = parseInt(this.cautionValue,10)
      let w = parseInt(this.warningValue,10)
      if (v < c) {
        return 0 < this.strokeColors.length ? this.strokeColors[0] : '#7fff7f';
      }else if (v < w) {
        return 1 < this.strokeColors.length ? this.strokeColors[1] : '#ffb266';
      }
      return 2 < this.strokeColors.length ? this.strokeColors[2] : '#ff7f7f';
    },
    strokeBgColor() {
      return 0 < this.strokeBgColors.length ? this.strokeBgColors[0] : '#ccffcc';
    },
    strokeCautionBgColor() {
      return 1 < this.strokeBgColors.length ? this.strokeBgColors[1] : '#ffd8b2';
    },
    strokeWarningBgColor() {
      return 2 < this.strokeBgColors.length ? this.strokeBgColors[2] : '#ffcccc';
    },
  }, 
  mounted() {
    this.startAnimate();
  },
  methods: {
    startAnimate: function() {
      this.isAnimating = true;
      this.speedEaseoutPercent = 0;
      this.animate();
    },
    animate: function() {
      if (isNaN(this.speedEaseoutPercent) && this.isAnimating) {
        this.isAnimating = false;
      }
      if (this.isAnimating) {
        this.speedEaseoutPercent += (this.percent - this.speedEaseoutPercent) * 0.1; //フレーム更新毎に加算
         
        if(this.speedEaseoutPercent > this.percent - this.percent * 0.005){
          this.isAnimating = false;
        }else {
          this.animation = requestAnimationFrame(this.animate);
        }
      }
    }
  },
  data() {
    return {
      animation: {},
      isAnimating: false,
      speedEaseoutPercent: 0,
    }
  }
}
</script>

<style lang="scss">
.guage {
  position: relative;
  float: auto;
  display: inline-block;
}

</style>


