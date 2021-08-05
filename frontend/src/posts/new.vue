<template>
  <Form class="col-sm-10" v-bind:obj="post" />
</template>

<script>
import Form from './form.vue'

export default {
  components: {
    Form
  },
  data () {
    return {
      post: {}
    }
  }
}
</script>

<style scoped>
</style>
