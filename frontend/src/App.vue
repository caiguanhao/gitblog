<template>
  <header>
    <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
      <div class="container">
        <router-link class="navbar-brand"
          v-bind:to="{ name: 'RouteHome' }">BLOG ADMIN</router-link>
        <button class="navbar-toggler" type="button"
          data-bs-toggle="collapse" data-bs-target="#navbarCollapse"
          aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarCollapse">
          <ul class="navbar-nav ms-auto mb-2 mb-md-0">
            <li class="nav-item">
              <router-link v-bind:to="{ name: 'RouteSync' }" class="nav-link">
                <strong class="me-1" v-text="currentStatus.UnpushedCommits"></strong>
                <small>commits unsynced</small>
              </router-link>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  </header>

  <main class="flex-shrink-0">
    <div class="container">
      <router-view></router-view>
    </div>
  </main>

  <footer class="footer mt-auto py-3 bg-light">
    <div class="container">
      <a class="text-muted text-decoration-none small"
        href="https://github.com/caiguanhao/gitblog"
        target="_blank">gitblog</a>
    </div>
  </footer>
</template>

<script>
export default {
  created () {
    document.addEventListener('click', (e) => {
      if (!e.isTrusted) return
      // don't click the row if text is selected before mouse button is released
      if (window.getSelection().toString().length) {
        return
      }
      let el = e.target
      while (el) {
        let node = el.nodeName
        if (node === 'A' || node === 'BUTTON' || node === 'INPUT') {
          if (node === 'A' && el.getAttribute('href')) {
            let path = el.pathname + el.search + el.hash
            if (this.$route.fullPath === path) { // "reload" if clicking the same route
              this.$router.replace({
                name: 'RouteBlank'
              }).then(() => {
                this.$router.replace(path)
              })
            }
            return
          }
          return
        }
        if (el.classList && el.classList.contains('clickable-row')) {
          let elem = el.querySelector('.clickable-row-target') ||
            el.querySelector('input[type=checkbox]') ||
            el.querySelector('a')
          if (elem) elem.click()
          return
        }
        el = el.parentNode
      }
      return
    })
  }
}
</script>

<style>
main > .container {
  padding-top: 80px;
}

.clickable-row {
  cursor: pointer;
}

.clickable-row:hover {
  background: #f1f1f1;
}

.fake-input { /* this prevents autofill */
  position: absolute;
  left: -100%;
  opacity: 0;
}
</style>
