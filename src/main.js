import { createApp } from 'vue'
import { Quasar } from 'quasar'
import quasarIconSet from 'quasar/icon-set/fontawesome-v6'
import '@quasar/extras/fontawesome-v6/fontawesome-v6.css'
import 'quasar/dist/quasar.css'
import App from './App.vue'
import router from './router'

const myApp = createApp(App)

myApp.use(Quasar, {
  iconSet: quasarIconSet,
  config: {
    dark: true,
  },
})

myApp.use(router)

myApp.mount('#app')
