
<script>
import LineChart from '@/components/modules/LineChart.vue';
import 'chartjs-plugin-streaming';

export default {
  name: 'dashboard',
  components: {
    LineChart
  },
  data: () => ({
    readDataset: null,
    writeDataset: null,
    readData: null,
    writeData: null,
    readOption: null,
    writeOption: null,
    socket: null
  }),
  mounted () {
    const socket = new WebSocket(`ws://${window.location.host}/ws`)
    this.$data.socket = socket
    this.$data.readData = new Map([])
    this.$data.writeData = new Map([])
    socket.onmessage = e => {
      const data = JSON.parse(e.data)
      data.forEach((value) => {
        this.createReadData(value)
        this.createWriteData(value)
      })
      this.createReadDataset(data)
      this.createWriteDataset(data)
    }
  },
  methods: {
    createReadData(value) {
      if (this.readData[value.file_name]) {
        this.readData[value.file_name].push(value.read_bytes_per_sec)
      }else {
        this.readData[value.file_name] = [value.read_bytes_per_sec]
      }
    },
    createWriteData(data) {
      if (this.writeData[value.file_name]) {
        this.writeData[value.file_name].push(value.write_bytes_per_sec)
      }else {
        this.writeData[value.file_name] = [value.write_bytes_per_sec]
      }
    },
    createReadDataset (data) {
      if (this.readDataset === null) {
        this.readDataset = {
          labels: [],
          datasets: []
        }
      }
      for (let i = 0;i < this.readDataset.datasets.Length; i++ ) {
        if (this.readDataset.datasets[i].label == key) {
          readDataset.datasets[i].data = value;
          break;
        }
      }
      this.readDataset = 
    },
    createWriteDataset (data) {
    }
  },
  destoryed () {
    this.$data.socket.close()
  }
}
</script>
