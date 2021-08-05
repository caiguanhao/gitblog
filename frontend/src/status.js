import http from './http'
import { reactive } from 'vue'

const currentStatus = reactive({})

export default {
  computed: {
    currentStatus () {
      return currentStatus
    }
  },
  methods: {
    getCurrentStatus () {
      return http.get('/status').then(res => {
        for (let key in currentStatus) {
          delete(currentStatus[key])
        }
        for (let key in res.data) {
          currentStatus[key] = res.data[key]
        }
      }, () => {
        for (let key in currentStatus) {
          delete(currentStatus[key])
        }
      })
    }
  }
}
