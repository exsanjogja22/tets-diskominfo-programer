import './style.css'
import { fetchProducts, fetchProductById, fetchCategories, fetchSuppliers } from './api.js'

const app = document.getElementById('app')

const SORT_OPTIONS = [
  { value: 'product_name:asc', label: 'Name (A-Z)' },
  { value: 'product_name:desc', label: 'Name (Z-A)' },
  { value: 'unit_price:asc', label: 'Price (Low-High)' },
  { value: 'unit_price:desc', label: 'Price (High-Low)' },
  { value: 'units_in_stock:asc', label: 'Stock (Low-High)' },
  { value: 'units_in_stock:desc', label: 'Stock (High-Low)' },
]

let state = {
  categories: [],
  suppliers: [],
  products: [],
  meta: null,
  loading: false,
  error: null,
  filter: {
    category_id: '',
    supplier_id: '',
    min_price: '',
    max_price: '',
    keyword: '',
    sort: 'product_name:asc',
    page: 1,
    limit: 10,
  },
  detailProduct: null,
  detailLoading: false,
  detailError: null,
}

function buildQuery() {
  const f = state.filter
  const q = {
    page: f.page,
    limit: f.limit,
    sort: f.sort,
  }
  if (f.category_id) q.category_id = f.category_id
  if (f.supplier_id) q.supplier_id = f.supplier_id
  if (f.min_price !== '') q.min_price = f.min_price
  if (f.max_price !== '') q.max_price = f.max_price
  if (f.keyword.trim()) q.keyword = f.keyword.trim()
  return q
}

function defaultFilter() {
  return {
    category_id: '',
    supplier_id: '',
    min_price: '',
    max_price: '',
    keyword: '',
    sort: 'product_name:asc',
    page: 1,
    limit: 10,
  }
}

async function loadOptions() {
  try {
    const [catRes, supRes] = await Promise.all([fetchCategories(), fetchSuppliers()])
    state.categories = catRes.data || []
    state.suppliers = supRes.data || []
  } catch (e) {
    state.error = 'Gagal memuat kategori/supplier: ' + e.message
  }
  render()
}

async function loadProducts() {
  state.loading = true
  state.error = null
  render()
  try {
    const res = await fetchProducts(buildQuery())
    state.products = res.data || []
    state.meta = res.meta || null
  } catch (e) {
    state.error = e.message || 'Gagal memuat produk'
    state.products = []
    state.meta = null
  }
  state.loading = false
  render()
}

async function applyFilter() {
  state.filter.page = 1
  await loadProducts()
}

function resetFilter() {
  state.filter = defaultFilter()
  loadProducts()
}

async function openDetail(id) {
  state.detailProduct = null
  state.detailError = null
  state.detailLoading = true
  render()
  try {
    const res = await fetchProductById(id)
    state.detailProduct = res.data
  } catch (e) {
    state.detailError = e.message || 'Gagal memuat detail produk'
  }
  state.detailLoading = false
  render()
}

function closeDetail() {
  state.detailProduct = null
  state.detailError = null
  render()
}

function goPage(page) {
  if (page < 1) return
  const totalPages = state.meta?.pagination?.total_pages || 1
  if (page > totalPages) return
  state.filter.page = page
  loadProducts()
}

function onLimitChange(limit) {
  state.filter.limit = Number(limit)
  state.filter.page = 1
  loadProducts()
}

function renderFilterPanel() {
  const f = state.filter
  const catOpts = state.categories.map(c => `<option value="${c.category_id}" ${f.category_id == c.category_id ? 'selected' : ''}>${escapeHtml(c.category_name)}</option>`).join('')
  const supOpts = state.suppliers.map(s => `<option value="${s.supplier_id}" ${f.supplier_id == s.supplier_id ? 'selected' : ''}>${escapeHtml(s.company_name)}</option>`).join('')
  const sortOpts = SORT_OPTIONS.map(o => `<option value="${o.value}" ${f.sort === o.value ? 'selected' : ''}>${o.label}</option>`).join('')

  return `
    <div class="bg-white rounded-xl shadow-sm border border-slate-200/80 p-5 mb-6">
      <h3 class="text-sm font-semibold text-slate-700 mb-4">Filter</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 xl:grid-cols-6 gap-4 items-end">
        <div>
          <label class="block text-xs font-medium text-slate-500 mb-1">Category</label>
          <select id="filter-category" class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
            <option value="">All</option>
            ${catOpts}
          </select>
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-500 mb-1">Supplier</label>
          <select id="filter-supplier" class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
            <option value="">All</option>
            ${supOpts}
          </select>
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-500 mb-1">Min Price</label>
          <input type="number" id="filter-min-price" step="0.01" min="0" placeholder="0" value="${escapeHtml(f.min_price)}" class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-500 mb-1">Max Price</label>
          <input type="number" id="filter-max-price" step="0.01" min="0" placeholder="—" value="${escapeHtml(f.max_price)}" class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
        </div>
        <div class="lg:col-span-2">
          <label class="block text-xs font-medium text-slate-500 mb-1">Search by name</label>
          <input type="text" id="filter-keyword" placeholder="Product name..." value="${escapeHtml(f.keyword)}" class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
        </div>
      </div>
      <div class="flex flex-wrap gap-3 mt-4 items-center">
        <div class="flex items-center gap-2">
          <label class="text-xs font-medium text-slate-500">Sort:</label>
          <select id="filter-sort" class="rounded-lg border border-slate-300 px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500">
            ${sortOpts}
          </select>
        </div>
        <button type="button" id="btn-apply-filter" class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition">Apply Filter</button>
        <button type="button" id="btn-reset-filter" class="px-4 py-2 bg-slate-200 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-300 transition">Reset Filter</button>
      </div>
    </div>
  `
}

function renderTable() {
  if (state.loading) {
    return `
      <div class="bg-white rounded-xl shadow-sm border border-slate-200/80 p-12 text-center">
        <div class="inline-block w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
        <p class="mt-3 text-slate-500">Memuat produk...</p>
      </div>
    `
  }
  if (state.error) {
    return `
      <div class="bg-white rounded-xl shadow-sm border border-slate-200/80 p-8 text-center">
        <p class="text-red-600 font-medium">${escapeHtml(state.error)}</p>
        <button type="button" id="btn-retry" class="mt-4 px-4 py-2 bg-slate-200 text-slate-700 rounded-lg hover:bg-slate-300">Coba lagi</button>
      </div>
    `
  }
  if (!state.products.length) {
    return `
      <div class="bg-white rounded-xl shadow-sm border border-slate-200/80 p-12 text-center">
        <p class="text-slate-500 text-lg">No products found.</p>
      </div>
    `
  }

  const rows = state.products.map(p => `
    <tr class="border-b border-slate-100 hover:bg-slate-50/80 transition">
      <td class="px-4 py-3 text-sm text-slate-800">${escapeHtml(p.product_name)}</td>
      <td class="px-4 py-3 text-sm text-slate-600">${escapeHtml(p.category_name || '—')}</td>
      <td class="px-4 py-3 text-sm text-slate-600">${escapeHtml(p.supplier_name || '—')}</td>
      <td class="px-4 py-3 text-sm text-slate-700">$${Number(p.unit_price).toFixed(2)}</td>
      <td class="px-4 py-3 text-sm text-slate-700">${p.units_in_stock}</td>
      <td class="px-4 py-3">
        <span class="inline-flex px-2 py-0.5 rounded text-xs font-medium ${p.discontinued ? 'bg-amber-100 text-amber-800' : 'bg-emerald-100 text-emerald-800'}">${p.discontinued ? 'Discontinued' : 'Active'}</span>
      </td>
      <td class="px-4 py-3">
        <button type="button" data-product-id="${p.product_id}" class="view-detail px-3 py-1.5 bg-blue-600 text-white text-xs font-medium rounded-lg hover:bg-blue-700 transition">View Details</button>
      </td>
    </tr>
  `).join('')

  return `
    <div class="bg-white rounded-xl shadow-sm border border-slate-200/80 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left">
          <thead class="bg-slate-50 border-b border-slate-200">
            <tr>
              <th class="px-4 py-3 text-xs font-semibold text-slate-600 uppercase tracking-wider">Product Name</th>
              <th class="px-4 py-3 text-xs font-semibold text-slate-600 uppercase tracking-wider">Category</th>
              <th class="px-4 py-3 text-xs font-semibold text-slate-600 uppercase tracking-wider">Supplier</th>
              <th class="px-4 py-3 text-xs font-semibold text-slate-600 uppercase tracking-wider">Unit Price</th>
              <th class="px-4 py-3 text-xs font-semibold text-slate-600 uppercase tracking-wider">Units In Stock</th>
              <th class="px-4 py-3 text-xs font-semibold text-slate-600 uppercase tracking-wider">Status</th>
              <th class="px-4 py-3 text-xs font-semibold text-slate-600 uppercase tracking-wider"></th>
            </tr>
          </thead>
          <tbody>${rows}</tbody>
        </table>
      </div>
    </div>
  `
}

function renderPagination() {
  if (state.loading || state.error || !state.meta) return ''
  const pag = state.meta.pagination || {}
  const total = pag.total || 0
  const page = pag.page || 1
  const limit = pag.limit || 10
  const totalPages = pag.total_pages || 1
  const from = total === 0 ? 0 : (page - 1) * limit + 1
  const to = Math.min(page * limit, total)

  const pages = []
  for (let i = 1; i <= totalPages; i++) {
    if (i === 1 || i === totalPages || (i >= page - 2 && i <= page + 2)) {
      pages.push(i)
    } else if (pages[pages.length - 1] !== '...') {
      pages.push('...')
    }
  }

  return `
    <div class="flex flex-wrap items-center justify-between gap-4 mt-4">
      <p class="text-sm text-slate-600">Showing ${from}-${to} of ${total} products</p>
      <div class="flex items-center gap-2">
        <label class="text-xs text-slate-500">Per page:</label>
        <select id="pagination-limit" class="rounded border border-slate-300 px-2 py-1 text-sm">
          <option value="10" ${limit === 10 ? 'selected' : ''}>10</option>
          <option value="20" ${limit === 20 ? 'selected' : ''}>20</option>
          <option value="50" ${limit === 50 ? 'selected' : ''}>50</option>
        </select>
      </div>
      <div class="flex items-center gap-1">
        <button type="button" id="btn-prev" ${page <= 1 ? 'disabled' : ''} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-50">Previous</button>
        ${pages.map(p => p === '...' ? '<span class="px-2 text-slate-400">…</span>' : `<button type="button" data-page="${p}" class="page-btn px-3 py-1.5 rounded-lg border text-sm font-medium ${p === page ? 'bg-blue-600 text-white border-blue-600' : 'border-slate-300 hover:bg-slate-50'}">${p}</button>`).join('')}
        <button type="button" id="btn-next" ${page >= totalPages ? 'disabled' : ''} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-50">Next</button>
      </div>
    </div>
  `
}

function renderDetailModal() {
  if (!state.detailProduct && !state.detailLoading && !state.detailError) return ''

  if (state.detailLoading) {
    return `
      <div id="modal-backdrop" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
        <div class="bg-white rounded-xl shadow-xl max-w-lg w-full p-8 text-center">
          <div class="inline-block w-10 h-10 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
          <p class="mt-4 text-slate-500">Memuat detail produk...</p>
        </div>
      </div>
    `
  }

  if (state.detailError) {
    return `
      <div id="modal-backdrop" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
        <div class="bg-white rounded-xl shadow-xl max-w-md w-full p-6 text-center">
          <p class="text-red-600 font-medium">${escapeHtml(state.detailError)}</p>
          <button type="button" id="btn-detail-close" class="mt-4 px-4 py-2 bg-slate-200 text-slate-700 rounded-lg hover:bg-slate-300">Tutup</button>
        </div>
      </div>
    `
  }

  const d = state.detailProduct
  return `
    <div id="modal-backdrop" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 overflow-y-auto">
      <div class="bg-white rounded-xl shadow-xl max-w-2xl w-full my-8">
        <div class="p-6 border-b border-slate-200 flex items-center justify-between">
          <h2 class="text-xl font-bold text-slate-800">Product Detail</h2>
          <button type="button" id="btn-detail-close" class="p-2 rounded-lg hover:bg-slate-100 text-slate-500">&times;</button>
        </div>
        <div class="p-6 space-y-4 text-sm">
          <div class="grid grid-cols-2 gap-3">
            <div class="text-slate-500">Product ID</div><div class="font-medium">${d.product_id}</div>
            <div class="text-slate-500">Product Name</div><div class="font-medium">${escapeHtml(d.product_name)}</div>
            <div class="text-slate-500">Category</div><div>${escapeHtml(d.category_name || '—')}</div>
            <div class="text-slate-500">Supplier</div><div>${escapeHtml(d.supplier_name || '—')}</div>
            <div class="text-slate-500">Quantity Per Unit</div><div>${escapeHtml(d.quantity_per_unit || '—')}</div>
            <div class="text-slate-500">Unit Price</div><div>$${Number(d.unit_price).toFixed(2)}</div>
            <div class="text-slate-500">Units In Stock</div><div>${d.units_in_stock}</div>
            <div class="text-slate-500">Units On Order</div><div>${d.units_on_order}</div>
            <div class="text-slate-500">Reorder Level</div><div>${d.reorder_level}</div>
            <div class="text-slate-500">Discontinued</div><div><span class="inline-flex px-2 py-0.5 rounded text-xs font-medium ${d.discontinued ? 'bg-amber-100 text-amber-800' : 'bg-emerald-100 text-emerald-800'}">${d.discontinued ? 'Yes' : 'No'}</span></div>
            <div class="text-slate-500">Total Terjual</div><div class="font-medium">${d.total_sold ?? 0}</div>
          </div>
        </div>
        <div class="p-6 border-t border-slate-200 flex gap-3">
          <button type="button" id="btn-back-to-list" class="px-4 py-2 bg-slate-200 text-slate-700 font-medium rounded-lg hover:bg-slate-300 transition">Back to List</button>
          <button type="button" id="btn-edit" class="px-4 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition">Edit</button>
        </div>
      </div>
    </div>
  `
}

function escapeHtml(s) {
  if (s == null) return ''
  const div = document.createElement('div')
  div.textContent = s
  return div.innerHTML
}

function render() {
  app.innerHTML = `
    <div class="min-h-screen bg-slate-100">
      <header class="bg-white border-b border-slate-200 shadow-sm">
        <div class="max-w-7xl mx-auto px-4 py-4">
          <h1 class="text-xl font-bold text-slate-800">Northwind Product Management</h1>
        </div>
      </header>
      <main class="max-w-7xl mx-auto px-4 py-6">
        ${renderFilterPanel()}
        ${renderTable()}
        ${renderPagination()}
      </main>
    </div>
    ${renderDetailModal()}
  `

  // Filter bindings
  const cat = document.getElementById('filter-category')
  const sup = document.getElementById('filter-supplier')
  const minP = document.getElementById('filter-min-price')
  const maxP = document.getElementById('filter-max-price')
  const kw = document.getElementById('filter-keyword')
  const sort = document.getElementById('filter-sort')
  if (cat) cat.addEventListener('change', () => { state.filter.category_id = cat.value; })
  if (sup) sup.addEventListener('change', () => { state.filter.supplier_id = sup.value; })
  if (minP) minP.addEventListener('input', () => { state.filter.min_price = minP.value; })
  if (maxP) maxP.addEventListener('input', () => { state.filter.max_price = maxP.value; })
  if (kw) kw.addEventListener('input', () => { state.filter.keyword = kw.value; })
  if (sort) sort.addEventListener('change', () => { state.filter.sort = sort.value; })

  document.getElementById('btn-apply-filter')?.addEventListener('click', applyFilter)
  document.getElementById('btn-reset-filter')?.addEventListener('click', resetFilter)
  document.getElementById('btn-retry')?.addEventListener('click', () => loadProducts())

  document.querySelectorAll('.view-detail').forEach(btn => {
    btn.addEventListener('click', () => openDetail(btn.dataset.productId))
  })

  document.getElementById('pagination-limit')?.addEventListener('change', e => onLimitChange(e.target.value))
  document.getElementById('btn-prev')?.addEventListener('click', () => goPage(state.filter.page - 1))
  document.getElementById('btn-next')?.addEventListener('click', () => goPage(state.filter.page + 1))
  document.querySelectorAll('.page-btn').forEach(btn => {
    btn.addEventListener('click', () => goPage(Number(btn.dataset.page)))
  })

  document.getElementById('btn-detail-close')?.addEventListener('click', closeDetail)
  document.getElementById('btn-back-to-list')?.addEventListener('click', closeDetail)
  document.getElementById('btn-edit')?.addEventListener('click', () => alert('Edit (UI only) - belum diimplementasi.'))
  document.getElementById('modal-backdrop')?.addEventListener('click', e => {
    if (e.target.id === 'modal-backdrop') closeDetail()
  })
}

async function init() {
  await loadOptions()
  await loadProducts()
}

init()
