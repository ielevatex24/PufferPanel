import Cookies from 'js-cookie'

export class InMemorySessionStore {
  _token = null
  _scopes = null
  _timeout = null

  setToken(token, expires) {
    clearTimeout(this._timeout)
    this._token = token
    setTimeout(() => {
      this._token = null
      this._scopes = null
    }, (expires * 1000) - Date.now())
  }

  setScopes(scopes) {
    this._scopes = scopes
  }

  getToken() {
    return this._token
  }

  getScopes() {
    return this._scopes
  }

  deleteSession() {
    this._token = null
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

  setToken(token, expires) {
    Cookies.set(AUTH_COOKIE_NAME, token, { ...this._cookieOptions, expires: new Date(expires * 1000) })
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
