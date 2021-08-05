<template>
  <form v-on:submit.prevent="submit">
    <div class="mb-3 row">
      <label class="col-sm-2 col-form-label">Title</label>
      <div class="col-sm-8">
        <input type="text" class="form-control" v-model="obj.Title" ref="Title">
      </div>
    </div>
    <div class="mb-3 row">
      <label class="col-sm-2 col-form-label">Body</label>
      <div class="col-sm-10">
        <textarea type="text" class="form-control" rows="15"
          v-model="obj.Body" ref="Body"></textarea>
      </div>
    </div>
    <div class="mb-3 row">
      <div class="col-sm-10 offset-sm-2">
        <button type="submit" class="btn btn-primary" v-bind:disabled="loading">Submit</button>
      </div>
    </div>
  </form>
</template>

<script>
import http from '../http'

export default {
  props: {
    obj: Object
  },
  data () {
    return {
      loading: false
    }
  },
  methods: {
    submit () {
      this.loading = true
      this.processErrors()
      if (this.obj.Id) {
        http.put(`/posts/${this.obj.Id}`, this.obj).then(res => {
          this.loading = false
          for (let key in res.data) {
            this.obj[key] = res.data[key]
          }
          this.$router.replace({
            name: 'RouteBlank'
          }).then(() => {
            this.$router.replace({
              name: 'RoutePostsEdit',
              params: {
                id: res.data.Id
              }
            })
          })
          this.$toast().success('Successfully updated post')
        }, (e) => {
          this.loading = false
          if (!this.processErrors(e)) {
            if (!e || !e.toastShown) {
              this.$toast().error('Error updating post')
            }
          }
        })
        return
      }
      http.post(`/posts`, this.obj).then(res => {
        this.loading = false
        this.$router.push({
          name: 'RouteHome'
        })
        this.$toast().success('Successfully created post')
      }, (e) => {
        this.loading = false
        if (!this.processErrors(e)) {
          if (!e || !e.toastShown) {
            this.$toast().error('Error creating post')
          }
        }
      })
    }
  }
}
</script>

<style scoped>
</style>
