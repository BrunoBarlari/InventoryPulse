<script>
  import { onMount, onDestroy } from 'svelte';
  import { categories as categoriesAPI, products as productsAPI } from '../lib/api.js';
  import { notifications } from '../lib/stores/notifications.js';
  import { websocketStore } from '../lib/stores/websocket.js';
  import { isAdmin } from '../lib/stores/auth.js';
  import Navbar from '../lib/components/Navbar.svelte';
  import Modal from '../lib/components/Modal.svelte';

  // State
  let categoriesData = { data: [], total_items: 0 };
  let productsData = { data: [], total_items: 0 };
  let isLoading = true;
  let activeTab = 'products';

  // Modal state
  let showCategoryModal = false;
  let showProductModal = false;
  let editingCategory = null;
  let editingProduct = null;

  // Form data
  let categoryForm = { name: '', description: '' };
  let productForm = { name: '', description: '', sku: '', quantity: 0, price: 0, category_id: 0 };

  // WebSocket event handlers
  let unsubscribers = [];

  onMount(async () => {
    await loadData();
    websocketStore.connect();
    unsubscribers = [
      websocketStore.on('product.created', handleProductCreated),
      websocketStore.on('product.updated', handleProductUpdated),
      websocketStore.on('product.deleted', handleProductDeleted),
      websocketStore.on('stock.updated', handleStockUpdated),
    ];
  });

  onDestroy(() => {
    unsubscribers.forEach(unsub => unsub && unsub());
  });

  function handleProductCreated(product) {
    productsData = {
      ...productsData,
      data: [...productsData.data, product],
      total_items: productsData.total_items + 1,
    };
  }

  function handleProductUpdated(product) {
    productsData = {
      ...productsData,
      data: productsData.data.map(p => p.id === product.id ? product : p),
    };
  }

  function handleProductDeleted(payload) {
    productsData = {
      ...productsData,
      data: productsData.data.filter(p => p.id !== payload.id),
      total_items: productsData.total_items - 1,
    };
  }

  function handleStockUpdated(product) {
    handleProductUpdated(product);
  }

  async function loadData() {
    isLoading = true;
    try {
      const [cats, prods] = await Promise.all([
        categoriesAPI.list(1, 100),
        productsAPI.list(1, 100),
      ]);
      categoriesData = cats;
      productsData = prods;
    } catch (err) {
      notifications.error('Failed to load data');
    } finally {
      isLoading = false;
    }
  }

  // Category CRUD
  function openCategoryModal(category = null) {
    editingCategory = category;
    categoryForm = category
      ? { name: category.name, description: category.description }
      : { name: '', description: '' };
    showCategoryModal = true;
  }

  async function saveCategory() {
    try {
      if (editingCategory) {
        await categoriesAPI.update(editingCategory.id, categoryForm);
        notifications.success('Category updated');
      } else {
        await categoriesAPI.create(categoryForm);
        notifications.success('Category created');
      }
      showCategoryModal = false;
      await loadData();
    } catch (err) {
      notifications.error(err.message || 'Failed to save category');
    }
  }

  async function deleteCategory(id) {
    if (!confirm('Are you sure you want to delete this category?')) return;
    try {
      await categoriesAPI.delete(id);
      notifications.success('Category deleted');
      await loadData();
    } catch (err) {
      notifications.error(err.message || 'Cannot delete category with products');
    }
  }

  // Product CRUD
  function openProductModal(product = null) {
    editingProduct = product;
    productForm = product
      ? { ...product }
      : { name: '', description: '', sku: '', quantity: 0, price: 0, category_id: categoriesData.data[0]?.id || 0 };
    showProductModal = true;
  }

  async function saveProduct() {
    try {
      const data = {
        ...productForm,
        price: parseFloat(productForm.price),
        quantity: parseInt(productForm.quantity),
        category_id: parseInt(productForm.category_id),
      };

      if (editingProduct) {
        await productsAPI.update(editingProduct.id, data);
        notifications.success('Product updated');
      } else {
        await productsAPI.create(data);
        notifications.success('Product created');
      }
      showProductModal = false;
    } catch (err) {
      notifications.error(err.message || 'Failed to save product');
    }
  }

  async function deleteProduct(id) {
    if (!confirm('Are you sure you want to delete this product?')) return;
    try {
      await productsAPI.delete(id);
      notifications.success('Product deleted');
    } catch (err) {
      notifications.error(err.message || 'Failed to delete product');
    }
  }

  async function updateStock(product, delta) {
    try {
      await productsAPI.updateStock(product.id, product.quantity + delta);
    } catch (err) {
      notifications.error('Failed to update stock');
    }
  }

  function getCategoryName(id) {
    const cat = categoriesData.data.find(c => c.id === id);
    return cat?.name || 'Unknown';
  }

  // Computed values
  $: totalStock = productsData.data.reduce((sum, p) => sum + p.quantity, 0);
  $: inventoryValue = productsData.data.reduce((sum, p) => sum + (p.price * p.quantity), 0);
  $: lowStockCount = productsData.data.filter(p => p.quantity < 10).length;
  $: avgPrice = productsData.data.length > 0
    ? productsData.data.reduce((sum, p) => sum + p.price, 0) / productsData.data.length
    : 0;

  // Category distribution for chart
  $: categoryDistribution = categoriesData.data.map(cat => {
    const products = productsData.data.filter(p => p.category_id === cat.id);
    const value = products.reduce((sum, p) => sum + (p.price * p.quantity), 0);
    return { name: cat.name, count: products.length, value };
  }).sort((a, b) => b.value - a.value);

  // Stock levels for chart (limit to 6 for better display)
  $: stockLevels = productsData.data
    .map(p => ({ name: p.name, quantity: p.quantity, low: p.quantity < 10 }))
    .sort((a, b) => b.quantity - a.quantity)
    .slice(0, 6);

  $: maxStock = Math.max(...stockLevels.map(s => s.quantity), 1);
  $: maxCategoryValue = Math.max(...categoryDistribution.map(c => c.value), 1);
</script>

<Navbar />

<main class="dashboard">
  <!-- Header with Stats -->
  <header class="dashboard-header">
    <div class="header-left">
      <h1 class="animate-title">Dashboard</h1>
      <p class="header-subtitle">Real-time inventory overview</p>
    </div>
    <div class="header-stats">
      <div class="mini-stat animate-stat" style="--delay: 0.1s">
        <span class="mini-stat-value">{productsData.total_items}</span>
        <span class="mini-stat-label">Products</span>
      </div>
      <div class="mini-stat animate-stat" style="--delay: 0.15s">
        <span class="mini-stat-value">{categoriesData.total_items}</span>
        <span class="mini-stat-label">Categories</span>
      </div>
      <div class="mini-stat animate-stat" style="--delay: 0.2s">
        <span class="mini-stat-value">{totalStock.toLocaleString()}</span>
        <span class="mini-stat-label">Total Stock</span>
      </div>
      <div class="mini-stat accent animate-stat" style="--delay: 0.25s">
        <span class="mini-stat-value">${(inventoryValue / 1000).toFixed(1)}k</span>
        <span class="mini-stat-label">Value</span>
      </div>
    </div>
  </header>

  {#if isLoading}
    <div class="loading-state">
      <div class="spinner spinner-lg"></div>
      <p>Loading data...</p>
    </div>
  {:else}
    <div class="dashboard-grid">
      <!-- Left: Data Tables -->
      <div class="main-content">
        <!-- Tabs -->
        <div class="tabs-container animate-slide-up" style="--delay: 0.2s">
          <div class="tabs">
            <button
              class="tab {activeTab === 'products' ? 'active' : ''}"
              on:click={() => activeTab = 'products'}
            >
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
              </svg>
              Products
              {#if lowStockCount > 0}
                <span class="tab-badge pulse">{lowStockCount}</span>
              {/if}
            </button>
            <button
              class="tab {activeTab === 'categories' ? 'active' : ''}"
              on:click={() => activeTab = 'categories'}
            >
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
              </svg>
              Categories
            </button>
          </div>
          {#if $isAdmin}
            <button class="btn btn-primary" on:click={() => activeTab === 'products' ? openProductModal() : openCategoryModal()}>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              Add {activeTab === 'products' ? 'Product' : 'Category'}
            </button>
          {/if}
        </div>

        <!-- Content -->
        {#if activeTab === 'products'}
          <div class="content-card card animate-slide-up" style="--delay: 0.3s">
            {#if productsData.data.length === 0}
              <div class="empty-state">
                <div class="empty-icon bounce">
                  <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                    <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                  </svg>
                </div>
                <h3>No products yet</h3>
                <p class="text-muted">Start by adding your first product</p>
              </div>
            {:else}
              <div class="table-container">
                <table>
                  <thead>
                    <tr>
                      <th>Product</th>
                      <th>SKU</th>
                      <th>Stock</th>
                      <th>Price</th>
                      {#if $isAdmin}<th class="text-right">Actions</th>{/if}
                    </tr>
                  </thead>
                  <tbody>
                    {#each productsData.data as product, i (product.id)}
                      <tr class="table-row-animate" style="--row-delay: {i * 0.05}s">
                        <td>
                          <div class="product-cell">
                            <strong>{product.name}</strong>
                            <span class="badge badge-neutral">{getCategoryName(product.category_id)}</span>
                          </div>
                        </td>
                        <td><code>{product.sku}</code></td>
                        <td>
                          <div class="stock-control">
                            {#if $isAdmin}
                              <button class="icon-btn" on:click={() => updateStock(product, -1)} disabled={product.quantity <= 0}>
                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                                  <line x1="5" y1="12" x2="19" y2="12"/>
                                </svg>
                              </button>
                            {/if}
                            <span class="stock-value {product.quantity < 10 ? 'low' : ''}">
                              {product.quantity}
                            </span>
                            {#if $isAdmin}
                              <button class="icon-btn" on:click={() => updateStock(product, 1)}>
                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                                  <line x1="12" y1="5" x2="12" y2="19"/>
                                  <line x1="5" y1="12" x2="19" y2="12"/>
                                </svg>
                              </button>
                            {/if}
                          </div>
                        </td>
                        <td class="price-cell">${product.price.toFixed(2)}</td>
                        {#if $isAdmin}
                          <td>
                            <div class="action-buttons">
                              <button class="btn btn-ghost btn-sm" on:click={() => openProductModal(product)}>Edit</button>
                              <button class="btn btn-danger btn-sm" on:click={() => deleteProduct(product.id)}>Delete</button>
                            </div>
                          </td>
                        {/if}
                      </tr>
                    {/each}
                  </tbody>
                </table>
              </div>
            {/if}
          </div>
        {:else}
          <div class="content-card card animate-slide-up" style="--delay: 0.3s">
            {#if categoriesData.data.length === 0}
              <div class="empty-state">
                <div class="empty-icon bounce">
                  <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                    <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
                  </svg>
                </div>
                <h3>No categories yet</h3>
                <p class="text-muted">Start by adding your first category</p>
              </div>
            {:else}
              <div class="categories-list">
                {#each categoriesData.data as category, i (category.id)}
                  <div class="category-row table-row-animate" style="--row-delay: {i * 0.05}s">
                    <div class="category-icon">
                      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
                      </svg>
                    </div>
                    <div class="category-info">
                      <h4>{category.name}</h4>
                      <p class="text-muted">{category.description || 'No description'}</p>
                    </div>
                    <div class="category-stats">
                      <span class="product-count">{productsData.data.filter(p => p.category_id === category.id).length} products</span>
                    </div>
                    {#if $isAdmin}
                      <div class="category-actions">
                        <button class="btn btn-ghost btn-sm" on:click={() => openCategoryModal(category)}>Edit</button>
                        <button class="btn btn-danger btn-sm" on:click={() => deleteCategory(category.id)}>Delete</button>
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}
      </div>

      <!-- Right: Analytics Panel -->
      <aside class="analytics-panel">
        <!-- Value by Category Chart -->
        <div class="chart-card card animate-slide-left" style="--delay: 0.4s">
          <div class="chart-header">
            <h3>Value by Category</h3>
            <span class="chart-total">${inventoryValue.toLocaleString()}</span>
          </div>
          <div class="bar-chart">
            {#if categoryDistribution.length === 0}
              <p class="text-muted text-center">No data yet</p>
            {:else}
              {#each categoryDistribution as cat, i}
                <div class="bar-row" style="--bar-delay: {0.5 + i * 0.1}s">
                  <div class="bar-label" title={cat.name}>{cat.name}</div>
                  <div class="bar-track">
                    <div
                      class="bar-fill animate-bar"
                      style="--width: {(cat.value / maxCategoryValue) * 100}%; --hue: {140 + i * 30}"
                    ></div>
                  </div>
                  <div class="bar-value">${(cat.value / 1000).toFixed(1)}k</div>
                </div>
              {/each}
            {/if}
          </div>
        </div>

        <!-- Stock Levels Chart -->
        <div class="chart-card card animate-slide-left" style="--delay: 0.5s">
          <div class="chart-header">
            <h3>Stock Levels</h3>
            {#if lowStockCount > 0}
              <span class="warning-badge pulse">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="12" y1="8" x2="12" y2="12"/>
                  <line x1="12" y1="16" x2="12.01" y2="16"/>
                </svg>
                {lowStockCount} low
              </span>
            {/if}
          </div>
          <div class="stock-chart-container">
            {#if stockLevels.length === 0}
              <p class="text-muted text-center">No products yet</p>
            {:else}
              <div class="stock-chart">
                {#each stockLevels as item, i}
                  <div class="stock-bar-wrapper" style="--bar-delay: {0.6 + i * 0.08}s">
                    <span class="stock-bar-value">{item.quantity}</span>
                    <div
                      class="stock-bar animate-grow"
                      class:low={item.low}
                      style="--height: {Math.max((item.quantity / maxStock) * 100, 8)}%"
                    ></div>
                    <span class="stock-bar-label" title={item.name}>{item.name.substring(0, 6)}{item.name.length > 6 ? '...' : ''}</span>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        </div>

        <!-- Quick Stats -->
        <div class="quick-stats card animate-slide-left" style="--delay: 0.6s">
          <div class="quick-stat">
            <div class="quick-stat-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div class="quick-stat-info">
              <span class="quick-stat-value">${avgPrice.toFixed(2)}</span>
              <span class="quick-stat-label">Avg. Price</span>
            </div>
          </div>
          <div class="divider"></div>
          <div class="quick-stat">
            <div class="quick-stat-icon warning">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
            </div>
            <div class="quick-stat-info">
              <span class="quick-stat-value">{lowStockCount}</span>
              <span class="quick-stat-label">Low Stock Items</span>
            </div>
          </div>
        </div>
      </aside>
    </div>
  {/if}
</main>

<!-- Category Modal -->
<Modal title={editingCategory ? 'Edit Category' : 'New Category'} show={showCategoryModal} on:close={() => showCategoryModal = false}>
  <form on:submit|preventDefault={saveCategory}>
    <div class="input-group">
      <label for="cat-name">Name</label>
      <input type="text" id="cat-name" class="input" bind:value={categoryForm.name} placeholder="Enter category name" required />
    </div>
    <div class="input-group">
      <label for="cat-desc">Description</label>
      <input type="text" id="cat-desc" class="input" bind:value={categoryForm.description} placeholder="Optional description" />
    </div>
    <div class="modal-actions">
      <button type="button" class="btn btn-secondary" on:click={() => showCategoryModal = false}>Cancel</button>
      <button type="submit" class="btn btn-primary">{editingCategory ? 'Update' : 'Create'}</button>
    </div>
  </form>
</Modal>

<!-- Product Modal -->
<Modal title={editingProduct ? 'Edit Product' : 'New Product'} show={showProductModal} on:close={() => showProductModal = false}>
  <form on:submit|preventDefault={saveProduct}>
    <div class="input-group">
      <label for="prod-name">Name</label>
      <input type="text" id="prod-name" class="input" bind:value={productForm.name} placeholder="Product name" required />
    </div>
    <div class="input-group">
      <label for="prod-sku">SKU</label>
      <input type="text" id="prod-sku" class="input" bind:value={productForm.sku} placeholder="SKU-001" required />
    </div>
    <div class="input-group">
      <label for="prod-desc">Description</label>
      <input type="text" id="prod-desc" class="input" bind:value={productForm.description} placeholder="Optional description" />
    </div>
    <div class="input-group">
      <label for="prod-cat">Category</label>
      <select id="prod-cat" class="input" bind:value={productForm.category_id}>
        {#each categoriesData.data as cat}
          <option value={cat.id}>{cat.name}</option>
        {/each}
      </select>
    </div>
    <div class="form-row">
      <div class="input-group">
        <label for="prod-qty">Quantity</label>
        <input type="number" id="prod-qty" class="input" bind:value={productForm.quantity} min="0" required />
      </div>
      <div class="input-group">
        <label for="prod-price">Price ($)</label>
        <input type="number" id="prod-price" class="input" bind:value={productForm.price} min="0.01" step="0.01" required />
      </div>
    </div>
    <div class="modal-actions">
      <button type="button" class="btn btn-secondary" on:click={() => showProductModal = false}>Cancel</button>
      <button type="submit" class="btn btn-primary">{editingProduct ? 'Update' : 'Create'}</button>
    </div>
  </form>
</Modal>

<style>
  .dashboard {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 24px 48px;
  }

  /* ===================== */
  /* HEADER */
  /* ===================== */
  .dashboard-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 32px 0 24px;
    gap: 24px;
    flex-wrap: wrap;
  }

  .header-left h1 {
    font-size: 2rem;
    margin-bottom: 4px;
  }

  .header-subtitle {
    color: var(--text-muted);
    font-size: 15px;
  }

  .header-stats {
    display: flex;
    gap: 12px;
  }

  .mini-stat {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 14px 24px;
    text-align: center;
    min-width: 100px;
  }

  .mini-stat.accent {
    background: var(--accent-light);
    border-color: transparent;
  }

  .mini-stat.accent .mini-stat-value {
    color: var(--accent-primary);
  }

  .mini-stat-value {
    display: block;
    font-family: var(--font-display);
    font-size: 1.5rem;
    line-height: 1.2;
    margin-bottom: 2px;
  }

  .mini-stat-label {
    font-size: 11px;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  /* ===================== */
  /* GRID LAYOUT */
  /* ===================== */
  .dashboard-grid {
    display: grid;
    grid-template-columns: 1fr 320px;
    gap: 24px;
    align-items: start;
  }

  @media (max-width: 1100px) {
    .dashboard-grid {
      grid-template-columns: 1fr;
    }
    .analytics-panel {
      order: -1;
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 16px;
    }
  }

  @media (max-width: 640px) {
    .dashboard-header {
      flex-direction: column;
      align-items: flex-start;
    }
    .header-stats {
      width: 100%;
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 8px;
    }
    .mini-stat {
      padding: 12px 16px;
      min-width: auto;
    }
    .mini-stat-value {
      font-size: 1.25rem;
    }
    .dashboard {
      padding: 0 16px 32px;
    }
  }

  /* ===================== */
  /* TABS */
  /* ===================== */
  .tabs-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    margin-bottom: 16px;
    flex-wrap: wrap;
  }

  .tabs {
    display: flex;
    gap: 4px;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 4px;
  }

  .tab {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: transparent;
    border: none;
    color: var(--text-muted);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    border-radius: var(--radius-md);
    transition: all 0.2s ease;
  }

  .tab:hover {
    color: var(--text-primary);
  }

  .tab.active {
    background: var(--bg-primary);
    color: var(--text-primary);
    box-shadow: var(--shadow-sm);
  }

  .tab-badge {
    background: var(--warning);
    color: white;
    font-size: 11px;
    font-weight: 600;
    padding: 2px 6px;
    border-radius: 100px;
  }

  /* ===================== */
  /* CONTENT CARD */
  /* ===================== */
  .content-card {
    padding: 0;
    overflow: hidden;
  }

  .loading-state, .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 24px;
    gap: 12px;
    text-align: center;
  }

  .loading-state {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    color: var(--text-muted);
  }

  .empty-icon {
    color: var(--text-muted);
    margin-bottom: 8px;
  }

  .empty-state h3 {
    font-size: 1rem;
  }

  /* ===================== */
  /* TABLE */
  /* ===================== */
  .product-cell {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .product-cell strong {
    font-weight: 500;
  }

  .product-cell .badge {
    width: fit-content;
    font-size: 11px;
    padding: 2px 8px;
  }

  .price-cell {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 14px;
  }

  .stock-control {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .stock-value {
    min-width: 36px;
    text-align: center;
    font-weight: 600;
    font-size: 14px;
    padding: 4px 8px;
    background: var(--bg-tertiary);
    border-radius: var(--radius-sm);
  }

  .stock-value.low {
    background: #FFFBEB;
    color: var(--warning);
  }

  .action-buttons {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }

  /* ===================== */
  /* CATEGORIES LIST */
  /* ===================== */
  .categories-list {
    display: flex;
    flex-direction: column;
  }

  .category-row {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border-color);
    transition: background 0.15s ease;
  }

  .category-row:last-child {
    border-bottom: none;
  }

  .category-row:hover {
    background: var(--bg-tertiary);
  }

  .category-icon {
    width: 40px;
    height: 40px;
    background: var(--accent-light);
    border-radius: var(--radius-md);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--accent-primary);
    flex-shrink: 0;
  }

  .category-info {
    flex: 1;
    min-width: 0;
  }

  .category-info h4 {
    font-size: 0.95rem;
    margin-bottom: 2px;
  }

  .category-info p {
    font-size: 13px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .category-stats {
    flex-shrink: 0;
  }

  .product-count {
    font-size: 13px;
    color: var(--text-muted);
    white-space: nowrap;
  }

  .category-actions {
    display: flex;
    gap: 8px;
    flex-shrink: 0;
  }

  /* ===================== */
  /* ANALYTICS PANEL */
  /* ===================== */
  .analytics-panel {
    display: flex;
    flex-direction: column;
    gap: 16px;
    position: sticky;
    top: 88px;
  }

  .chart-card {
    padding: 20px;
    overflow: hidden;
  }

  .chart-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
    gap: 12px;
  }

  .chart-header h3 {
    font-size: 0.9rem;
    font-family: var(--font-body);
    font-weight: 600;
    margin: 0;
  }

  .chart-total {
    font-family: var(--font-display);
    font-size: 1rem;
    color: var(--accent-primary);
    flex-shrink: 0;
  }

  .warning-badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background: #FFFBEB;
    color: var(--warning);
    font-size: 12px;
    font-weight: 500;
    padding: 4px 10px;
    border-radius: 100px;
    flex-shrink: 0;
  }

  /* ===================== */
  /* BAR CHART (Horizontal) */
  /* ===================== */
  .bar-chart {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .bar-row {
    display: grid;
    grid-template-columns: 60px 1fr 45px;
    align-items: center;
    gap: 10px;
  }

  .bar-label {
    font-size: 12px;
    color: var(--text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .bar-track {
    height: 8px;
    background: var(--bg-tertiary);
    border-radius: 4px;
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    background: hsl(var(--hue, 160), 60%, 40%);
    border-radius: 4px;
    width: var(--width);
    max-width: 100%;
  }

  .bar-value {
    font-size: 12px;
    font-weight: 500;
    text-align: right;
    color: var(--text-primary);
    flex-shrink: 0;
  }

  /* ===================== */
  /* STOCK CHART (Vertical Bars) */
  /* ===================== */
  .stock-chart-container {
    overflow: hidden;
  }

  .stock-chart {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    height: 120px;
    gap: 8px;
    padding-top: 24px;
  }

  .stock-bar-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 100%;
    min-width: 0;
    max-width: 40px;
  }

  .stock-bar {
    width: 100%;
    background: var(--accent-primary);
    border-radius: 4px 4px 0 0;
    height: var(--height);
    min-height: 4px;
    max-height: calc(100% - 24px);
    transition: background 0.2s ease;
  }

  .stock-bar.low {
    background: var(--warning);
  }

  .stock-bar-value {
    font-size: 10px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 4px;
    height: 14px;
  }

  .stock-bar-label {
    margin-top: 6px;
    font-size: 9px;
    color: var(--text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
    text-align: center;
  }

  /* ===================== */
  /* QUICK STATS */
  /* ===================== */
  .quick-stats {
    padding: 16px 20px;
  }

  .quick-stat {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 8px 0;
  }

  .quick-stat-icon {
    width: 40px;
    height: 40px;
    background: var(--accent-light);
    border-radius: var(--radius-md);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--accent-primary);
    flex-shrink: 0;
  }

  .quick-stat-icon.warning {
    background: #FFFBEB;
    color: var(--warning);
  }

  .quick-stat-info {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .quick-stat-value {
    font-family: var(--font-display);
    font-size: 1.25rem;
    line-height: 1.2;
  }

  .quick-stat-label {
    font-size: 12px;
    color: var(--text-muted);
  }

  .divider {
    height: 1px;
    background: var(--border-color);
    margin: 8px 0;
  }

  /* ===================== */
  /* MODAL */
  /* ===================== */
  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
    padding-top: 16px;
    border-top: 1px solid var(--border-color);
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
  }

  /* ===================== */
  /* ANIMATIONS */
  /* ===================== */
  .animate-title {
    animation: titleReveal 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
    opacity: 0;
  }

  @keyframes titleReveal {
    from {
      opacity: 0;
      transform: translateY(20px);
      filter: blur(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
      filter: blur(0);
    }
  }

  .animate-stat {
    animation: statPop 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
    animation-delay: var(--delay, 0s);
    opacity: 0;
    transform: scale(0.8);
  }

  @keyframes statPop {
    from {
      opacity: 0;
      transform: scale(0.8) translateY(10px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }

  .animate-slide-up {
    animation: slideUp 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
    animation-delay: var(--delay, 0s);
    opacity: 0;
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-slide-left {
    animation: slideLeft 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
    animation-delay: var(--delay, 0s);
    opacity: 0;
  }

  @keyframes slideLeft {
    from {
      opacity: 0;
      transform: translateX(30px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .table-row-animate {
    animation: rowSlide 0.4s ease forwards;
    animation-delay: var(--row-delay, 0s);
    opacity: 0;
  }

  @keyframes rowSlide {
    from {
      opacity: 0;
      transform: translateX(-20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .animate-bar {
    animation: barGrow 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
    animation-delay: var(--bar-delay, 0s);
    transform-origin: left;
    transform: scaleX(0);
  }

  @keyframes barGrow {
    from { transform: scaleX(0); }
    to { transform: scaleX(1); }
  }

  .animate-grow {
    animation: growUp 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
    animation-delay: var(--bar-delay, 0s);
    transform-origin: bottom;
    transform: scaleY(0);
  }

  @keyframes growUp {
    from { transform: scaleY(0); }
    to { transform: scaleY(1); }
  }

  .bounce {
    animation: bounce 2s infinite;
  }

  @keyframes bounce {
    0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
    40% { transform: translateY(-10px); }
    60% { transform: translateY(-5px); }
  }

  .pulse {
    animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.6; }
  }

  .icon-btn {
    transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .icon-btn:hover:not(:disabled) {
    transform: scale(1.1);
  }

  .icon-btn:active:not(:disabled) {
    transform: scale(0.95);
  }

  .btn {
    transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .btn:hover:not(:disabled) {
    transform: translateY(-2px);
  }

  .btn:active:not(:disabled) {
    transform: translateY(0);
  }

  .card {
    transition: box-shadow 0.2s ease, transform 0.2s ease;
  }

  .chart-card:hover {
    box-shadow: var(--shadow-md);
  }
</style>
