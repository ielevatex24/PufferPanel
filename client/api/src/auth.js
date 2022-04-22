function is2xx(status) {
  return status >= 200 && status < 300
}

export class AuthApi {
  _api = null
  _sessionStore = null

  constructor(api, sessionStore) {
    this._api = api
    this._sessionStore = sessionStore
  }

  _handleLogin(token, scopes) {
    const payload = JSON.parse(atob(token.split('.')[1]))
    this._sessionStore.setToken(token, payload.exp)
    this._sessionStore.setScopes(scopes)
    return true
  }

  async login(email, password) {
    const res = await this._api.post('/auth/login', { email, password })
    if (res.data.otpNeeded) return 'otp'
    return this._handleLogin(res.data.session, res.data.scopes)
  }

  async loginOtp(token) {
    const res = await this._api.post('/auth/otp', { token })
    return this._handleLogin(res.data.session, res.data.scopes)
  }

  async register(username, email, password) {
    const res = await this._api.post('/auth/register', { username, email, password })
    return this._handleLogin(res.data.token, res.data.scopes)
  }

  async reauth() {
    const res = await this._api.post('/auth/reauth')
    if (is2xx(res.status)) {
      this._handleLogin(res.data.session, res.data.scopes)
    }
  }

  getToken() {
    return this._sessionStore.getToken()
  }

  isLoggedIn() {
    return this.getToken() !== null
  }

  hasScope(scope) {
    if (!this.isLoggedIn()) return false
    const scopes = this._sessionStore.getScopes()
    if (scopes !== null) {
      if (scopes.indexOf(scope) !== -1) return true
      return scopes.indexOf('servers.admin') !== -1
    }
    return false
  }

  logout() {
    this._sessionStore.deleteSession()
  }
}
