<template>
  <router-link v-bind:to="{ name: 'RoutePostsNew' }"
    class="btn btn-primary mb-3">Add a New Post</router-link>
  <template v-if="!posts || !posts.length">
    <h5 class="mb-0 text-muted">No posts have been created yet</h5>
  </template>
  <div v-else class="table-responsive">
    <table class="table">
      <thead>
        <tr>
          <th>TITLE</th>
          <th width="30%">UPDATED AT</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="post in posts" class="clickable-row">
          <td>
            <router-link v-bind:to="{ name: 'RoutePostsEdit', params: { id: post.Id } }"
              class="text-break" v-text="post.Title"></router-link>
          </td>
          <td>
            <div class="text-nowrap" v-text="timeago(post.UpdatedAt)"></div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import http from '../http'
import * as timeago from 'timeago.js'

export default {
  data () {
    return {
      posts: []
    }
  },
  methods: {
    timeago (time) {
      return timeago.format(time)
    }
  },
  beforeRouteEnter (to, from, next) {
    http.get('/posts', { params: to.query }).then(res => {
      next(vm => {
        vm.posts = res.data
      })
    }, next)
  },
  beforeRouteUpdate (to, from, next) {
    http.get('/posts', { params: to.query }).then(res => {
      this.posts = res.data
      next()
    }, next)
  }
}
</script>
