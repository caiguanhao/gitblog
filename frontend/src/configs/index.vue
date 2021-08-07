<template>
  <form class="col-sm-8 offset-sm-2" v-on:submit.prevent="submit">
    <input type="text" class="fake-input">
    <input type="password" class="fake-input">
    <div class="mb-3" v-for="config in configs">
      <label class="form-label" v-text="config.Key"></label>
      <textarea type="text" class="form-control" rows="6"
        v-if="types[config.Key] === 'textarea'" v-model="values[config.Key]"></textarea>
      <input type="password" class="form-control"
        v-else-if="types[config.Key] === 'password'" v-model="values[config.Key]">
      <input type="text" class="form-control"
        v-else v-model="values[config.Key]">
      <div class="form-text" v-text="config.Comment"></div>
    </div>
    <div class="mb-3">
      <button type="submit" class="btn btn-primary"
        v-bind:disabled="loading" v-text="buttonText"></button>
    </div>
  </form>
</template>

<script>
import http from '../http'

export default {
  data () {
    return {
      loading: false,
      restarting: false,
      configs: [],
      values: {},
      types: {
        SSHPrivateKey: 'textarea',
        SSHPrivateKeyPassword: 'password'
      }
    }
  },
  computed: {
    buttonText () {
      if (this.restarting) return 'Restarting server...'
      if (this.loading) return 'Loading...'
      return 'Submit'
    }
  },
  methods: {
    submit () {
      this.loading = true
      let data = JSON.parse(JSON.stringify(this.values))
      if (data.SSHPrivateKey.slice(-1) !== '\n') data.SSHPrivateKey += '\n'
      data.SSHPrivateKey = window.btoa(data.SSHPrivateKey)
      http.post(`/configs`, data).then(res => {
        this.restarting = true
        setTimeout(() => {
          this.loading = false
          this.restarting = false
          this.$router.push({
            name: 'RouteHome'
          }).then(() => {
            this.$toast().success('Successfully updated configs')
          })
        }, 3000)
      }, (e) => {
        this.loading = false
        if (!this.processErrors(e)) {
          if (!e || !e.toastShown) {
            this.$toast().error('Failed to update configs')
          }
        }
      })
    }
  },
  beforeRouteEnter (to, from, next) {
    http.get('/configs').then(res => {
      next(vm => {
        vm.configs = res.data
        vm.values = {}
        vm.configs.forEach(config => {
          vm.values[config.Key] = config.Value.trim()
        })
      })
    }, next)
  },
  beforeRouteUpdate (to, from, next) {
    http.get('/configs').then(res => {
      this.configs = res.data
      this.values = {}
      this.configs.forEach(config => {
        this.values[config.Key] = config.Value.trim()
      })
      next()
    }, next)
  }
}
</script>
