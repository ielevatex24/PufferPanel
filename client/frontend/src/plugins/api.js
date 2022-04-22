import ApiClient, { CookieSessionStore } from 'pufferpanel'

export const apiClient = new ApiClient(
  location.origin,
  new CookieSessionStore()
)

export default {
  install: (app) => {
    app.config.globalProperties.$api = apiClient
    app.provide('api', apiClient)
  }
}
