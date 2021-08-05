<template>
  <Form class="col-sm-10 mb-5" v-bind:obj="post" />
  <hr class="mb-5" />
  <div class="col-sm-10 mb-5">
    <div class="row">
      <div class="offset-sm-2">
        <button type="button" class="btn btn-danger" v-on:click.prevent="destroy">
          Delete
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import http from '../http'
import Form from './form.vue'

export default {
  components: {
    Form
  },
  data () {
    return {
      post: {}
    }
  },
  methods: {
    destroy () {
      if (!window.confirm('Permanently delete post?')) return
      http.delete(`/posts/${this.post.Id}`).then(res => {
        this.$toast().success('Successfully deleted post')
        this.$router.push({
          name: 'RouteHome'
        })
      }, (error) => {
        if (!error || !error.toastShown) {
          this.$toast().error('Error deleting post')
        }
      })
    }
  },
  beforeRouteEnter (to, from, next) {
    http.get(`/posts/${to.params.id}`).then(res => {
      next(vm => {
        vm.post = res.data
      })
    }, next)
  },
  beforeRouteUpdate (to, from, next) {
    http.get(`/posts/${to.params.id}`).then(res => {
      this.post = res.data
      next()
    }, next)
  }
}
</script>

<style scoped>
</style>
