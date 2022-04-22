export class TemplateApi {
  _api = null

  constructor(api) {
    this._api = api
  }

  async list() {
    const res = await this._api.get('/api/templates')
    return res.data
  }

  async get(name) {
    const res = await this._api.get(`/api/templates/${name}`)
    return res.data
  }

  async save(name, template) {
    await this._api.put(`/api/templates/${name}`, template)
    return true
  }

  async delete(name) {
    await this._api.delete(`/api/templates/${name}`)
    return true
  }

  async listImportable() {
    const res = await this._api.post('/api/templates/import')
    return res.data
  }

  async import(name) {
    await this._api.post(`/api/templates/import/${name}`)
    return true
  }
}
