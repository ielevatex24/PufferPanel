import Cookies from 'js-cookie'

export class InMemorySessionStore {
  _token = null
  _scopes = null
  _timeout = null

  setScopes(scopes) {
    this._scopes = scopes
  }

  getScopes() {
    return this._scopes
  }

  deleteSession() {
    this._scopes = null
  }
}

const defaultCookieSessionStoreOptions = {
  domain: undefined,
  path: '/',
  secure: window ? window.location.protocol === 'https' : false,
  sameSite: 'Strict'
}

const AUTH_COOKIE_NAME = 'puffer_auth'
const SCOPES_COOKIE_NAME = 'puffer_scopes'
export class CookieSessionStore {
  _cookieOptions = null

  constructor(options = defaultCookieSessionStoreOptions) {
    this._cookieOptions = options
  }

  setScopes(scopes) {
    Cookies.set(SCOPES_COOKIE_NAME, JSON.stringify(scopes), this._cookieOptions)
  }

  getToken() {
    return Cookies.get(AUTH_COOKIE_NAME) || null
  }

  getScopes() {
    const res = Cookies.get(SCOPES_COOKIE_NAME)
    if (res) return JSON.parse(res)
    return null
  }

  deleteSession() {
    Cookies.remove(AUTH_COOKIE_NAME, this._cookieOptions)
    Cookies.remove(SCOPES_COOKIE_NAME, this._cookieOptions)
  }
}
