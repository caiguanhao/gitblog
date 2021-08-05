<template>
  <div class="p-4">
    <div class="row">
      <div class="col-sm-6 offset-sm-3">
        <div class="card">
          <div class="card-body text-center">
            <template v-if="statusCode === 404">
              <h3>404</h3>
              <h5>The webpage you visited does not exist.</h5>
            </template>
            <template v-else-if="statusCode === 403">
              <h3>403</h3>
              <h5>You do not have permission to access this page.</h5>
            </template>
            <template v-else>
              <h3>500</h3>
              <h5>Server error, unable to process request. Please try again later.</h5>
            </template>
            <router-link v-bind:to="$router.$lastRoute" class="btn btn-primary mt-3"
              v-if="$router.$lastRoute" v-text="`Try again: ${$router.$lastRoute.fullPath}`"></router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import router from '../routes'

export default {
  computed: {
    statusCode () {
      let err = router.$lastError
      if (!err) return 404
      if (err && err.response) {
        return err.response.status
      }
      return null
    }
  },
  beforeRouteUpdate (to, from, next) {
    if (!router.$lastRoute) {
      next({ name: 'RouteHome' })
      return
    }
    next()
  },
  beforeRouteEnter (to, from, next) {
    if (!router.$lastRoute) {
      next({ name: 'RouteHome' })
      return
    }
    next()
  }
}
</script>

<style scoped>
h3 {
  font-size: 148px;
  margin: 0px;
  color: #ff508e;
}

.btn {
  max-width: 300px;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}
</style>
