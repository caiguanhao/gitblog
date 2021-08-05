<template>
  <div class="d-flex flex-column justify-content-center align-items-center box">
    <div class="mb-4">
      <small>You have</small>
      <strong class="mx-1" v-text="currentStatus.UnpushedCommits"></strong>
      <small>commits to sync.</small>
    </div>
    <button type="button" class="btn btn-primary"
      v-on:click.prevent="sync"
      v-text="loading ? 'Syncing...' : 'Sync Now'"
      v-bind:disabled="loading"></button>
  </div>
</template>

<script>
import http from '../http'

export default {
  data () {
    return {
      loading: false
    }
  },
  methods: {
    sync () {
      this.loading = true
      http.post('/push', null, {
        timeout: 20000
      }).then(() => {
        return this.getCurrentStatus()
      }).then(() => {
        this.loading = false
        this.$toast().success('Successfully synced')
      }, (err) => {
        console.error(err)
        this.loading = false
        this.$toast().error('Failed to sync')
      })
    }
  }
}
</script>

<style scoped>
.box {
  min-height: 300px;
}
</style>
