const API_BASE = ''

export async function fetchProducts(params = {}) {
  const q = new URLSearchParams()
  if (params.page) q.set('page', params.page)
  if (params.limit) q.set('limit', params.limit)
  if (params.category_id) q.set('category_id', params.category_id)
  if (params.supplier_id) q.set('supplier_id', params.supplier_id)
  if (params.min_price != null && params.min_price !== '') q.set('min_price', params.min_price)
  if (params.max_price != null && params.max_price !== '') q.set('max_price', params.max_price)
  if (params.keyword) q.set('keyword', params.keyword)
  if (params.sort) q.set('sort', params.sort)
  const res = await fetch(`${API_BASE}/api/products?${q}`)
  if (!res.ok) throw new Error(await res.text() || res.statusText)
  return res.json()
}

export async function fetchProductById(id) {
  const res = await fetch(`${API_BASE}/api/products/${id}`)
  if (!res.ok) throw new Error(await res.text() || res.statusText)
  return res.json()
}

export async function fetchCategories() {
  const res = await fetch(`${API_BASE}/api/categories`)
  if (!res.ok) throw new Error(await res.text() || res.statusText)
  return res.json()
}

export async function fetchSuppliers() {
  const res = await fetch(`${API_BASE}/api/suppliers`)
  if (!res.ok) throw new Error(await res.text() || res.statusText)
  return res.json()
}
