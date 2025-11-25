<script>
  import { onMount } from 'svelte';
  import { categories as categoriesAPI, products as productsAPI } from '../lib/api.js';
  import { notifications } from '../lib/stores/notifications.js';
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

  onMount(async () => {
    await loadData();
  });

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
      await loadData();
    } catch (err) {
      notifications.error(err.message || 'Failed to save product');
    }
  }

  async function deleteProduct(id) {
    if (!confirm('Are you sure you want to delete this product?')) return;
    try {
      await productsAPI.delete(id);
      notifications.success('Product deleted');
      await loadData();
    } catch (err) {
      notifications.error(err.message || 'Failed to delete product');
    }
  }

  async function updateStock(product, delta) {
    try {
      await productsAPI.updateStock(product.id, product.quantity + delta);
      notifications.success(`Stock updated: ${product.quantity + delta}`);
      await loadData();
    } catch (err) {
      notifications.error('Failed to update stock');
    }
  }

  function getCategoryName(id) {
    const cat = categoriesData.data.find(c => c.id === id);
    return cat?.name || 'Unknown';
  }
</script>

<Navbar />

<main class="dashboard">
  <div class="dashboard-header">
    <div class="stats-grid">
      <div class="stat-card glass">
        <div class="stat-icon">📦</div>
        <div class="stat-content">
          <span class="stat-value">{productsData.total_items}</span>
          <span class="stat-label">Products</span>
        </div>
      </div>
      <div class="stat-card glass">
        <div class="stat-icon">📁</div>
        <div class="stat-content">
          <span class="stat-value">{categoriesData.total_items}</span>
          <span class="stat-label">Categories</span>
        </div>
      </div>
      <div class="stat-card glass">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <span class="stat-value">
            {productsData.data.reduce((sum, p) => sum + p.quantity, 0)}
          </span>
          <span class="stat-label">Total Stock</span>
        </div>
      </div>
      <div class="stat-card glass">
        <div class="stat-icon">💰</div>
        <div class="stat-content">
          <span class="stat-value">
            ${productsData.data.reduce((sum, p) => sum + (p.price * p.quantity), 0).toLocaleString('en-US', { minimumFractionDigits: 0, maximumFractionDigits: 0 })}
          </span>
          <span class="stat-label">Inventory Value</span>
        </div>
      </div>
    </div>
  </div>

  <div class="tabs">
    <button 
      class="tab {activeTab === 'products' ? 'active' : ''}"
      on:click={() => activeTab = 'products'}
    >
      📦 Products
    </button>
    <button 
      class="tab {activeTab === 'categories' ? 'active' : ''}"
      on:click={() => activeTab = 'categories'}
    >
      📁 Categories
    </button>
  </div>

  {#if isLoading}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Loading data...</p>
    </div>
  {:else if activeTab === 'products'}
    <div class="content-card glass animate-fadeIn">
      <div class="content-header">
        <h2>Products</h2>
        {#if $isAdmin}
          <button class="btn btn-primary" on:click={() => openProductModal()}>
            + Add Product
          </button>
        {/if}
      </div>

      {#if productsData.data.length === 0}
        <div class="empty-state">
          <span class="empty-icon">📦</span>
          <p>No products yet</p>
        </div>
      {:else}
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>SKU</th>
                <th>Category</th>
                <th>Stock</th>
                <th>Price</th>
                {#if $isAdmin}<th>Actions</th>{/if}
              </tr>
            </thead>
            <tbody>
              {#each productsData.data as product}
                <tr>
                  <td>
                    <strong>{product.name}</strong>
                    {#if product.description}
                      <br /><small class="text-muted">{product.description}</small>
                    {/if}
                  </td>
                  <td><code>{product.sku}</code></td>
                  <td>
                    <span class="badge badge-info">{getCategoryName(product.category_id)}</span>
                  </td>
                  <td>
                    <div class="stock-control">
                      {#if $isAdmin}
                        <button class="btn btn-icon btn-secondary btn-sm" on:click={() => updateStock(product, -1)}>-</button>
                      {/if}
                      <span class="stock-value {product.quantity < 10 ? 'low' : ''}">
                        {product.quantity}
                      </span>
                      {#if $isAdmin}
                        <button class="btn btn-icon btn-secondary btn-sm" on:click={() => updateStock(product, 1)}>+</button>
                      {/if}
                    </div>
                  </td>
                  <td>${product.price.toFixed(2)}</td>
                  {#if $isAdmin}
                    <td>
                      <div class="action-buttons">
                        <button class="btn btn-secondary btn-sm" on:click={() => openProductModal(product)}>Edit</button>
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
    <div class="content-card glass animate-fadeIn">
      <div class="content-header">
        <h2>Categories</h2>
        {#if $isAdmin}
          <button class="btn btn-primary" on:click={() => openCategoryModal()}>
            + Add Category
          </button>
        {/if}
      </div>

      {#if categoriesData.data.length === 0}
        <div class="empty-state">
          <span class="empty-icon">📁</span>
          <p>No categories yet</p>
        </div>
      {:else}
        <div class="categories-grid">
          {#each categoriesData.data as category}
            <div class="category-card glass glass-hover">
              <div class="category-info">
                <h4>{category.name}</h4>
                <p class="text-muted">{category.description || 'No description'}</p>
                <span class="text-secondary">
                  {productsData.data.filter(p => p.category_id === category.id).length} products
                </span>
              </div>
              {#if $isAdmin}
                <div class="category-actions">
                  <button class="btn btn-secondary btn-sm" on:click={() => openCategoryModal(category)}>Edit</button>
                  <button class="btn btn-danger btn-sm" on:click={() => deleteCategory(category.id)}>Delete</button>
                </div>
              {/if}
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</main>

<!-- Category Modal -->
<Modal title={editingCategory ? 'Edit Category' : 'New Category'} show={showCategoryModal} on:close={() => showCategoryModal = false}>
  <form on:submit|preventDefault={saveCategory}>
    <div class="input-group">
      <label for="cat-name">Name</label>
      <input type="text" id="cat-name" class="input" bind:value={categoryForm.name} required />
    </div>
    <div class="input-group">
      <label for="cat-desc">Description</label>
      <input type="text" id="cat-desc" class="input" bind:value={categoryForm.description} />
    </div>
    <div class="flex gap-4 mt-4">
      <button type="button" class="btn btn-secondary" on:click={() => showCategoryModal = false}>Cancel</button>
      <button type="submit" class="btn btn-primary">Save</button>
    </div>
  </form>
</Modal>

<!-- Product Modal -->
<Modal title={editingProduct ? 'Edit Product' : 'New Product'} show={showProductModal} on:close={() => showProductModal = false}>
  <form on:submit|preventDefault={saveProduct}>
    <div class="input-group">
      <label for="prod-name">Name</label>
      <input type="text" id="prod-name" class="input" bind:value={productForm.name} required />
    </div>
    <div class="input-group">
      <label for="prod-sku">SKU</label>
      <input type="text" id="prod-sku" class="input" bind:value={productForm.sku} required />
    </div>
    <div class="input-group">
      <label for="prod-desc">Description</label>
      <input type="text" id="prod-desc" class="input" bind:value={productForm.description} />
    </div>
    <div class="input-group">
      <label for="prod-cat">Category</label>
      <select id="prod-cat" class="input" bind:value={productForm.category_id}>
        {#each categoriesData.data as cat}
          <option value={cat.id}>{cat.name}</option>
        {/each}
      </select>
    </div>
    <div class="flex gap-4">
      <div class="input-group" style="flex:1">
        <label for="prod-qty">Quantity</label>
        <input type="number" id="prod-qty" class="input" bind:value={productForm.quantity} min="0" required />
      </div>
      <div class="input-group" style="flex:1">
        <label for="prod-price">Price ($)</label>
        <input type="number" id="prod-price" class="input" bind:value={productForm.price} min="0.01" step="0.01" required />
      </div>
    </div>
    <div class="flex gap-4 mt-4">
      <button type="button" class="btn btn-secondary" on:click={() => showProductModal = false}>Cancel</button>
      <button type="submit" class="btn btn-primary">Save</button>
    </div>
  </form>
</Modal>

<style>
  .dashboard {
    padding: 0 32px 32px;
    max-width: 1400px;
    margin: 0 auto;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
    margin-bottom: 32px;
  }

  .stat-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 24px;
  }

  .stat-icon {
    font-size: 32px;
  }

  .stat-content {
    display: flex;
    flex-direction: column;
  }

  .stat-value {
    font-family: var(--font-display);
    font-size: 1.75rem;
    font-weight: 700;
  }

  .stat-label {
    color: var(--text-muted);
    font-size: 14px;
  }

  .tabs {
    display: flex;
    gap: 8px;
    margin-bottom: 24px;
  }

  .tab {
    padding: 12px 24px;
    background: transparent;
    border: 1px solid transparent;
    color: var(--text-muted);
    font-size: 15px;
    cursor: pointer;
    border-radius: var(--radius-md);
    transition: all 0.2s;
  }

  .tab:hover {
    color: var(--text-primary);
    background: var(--glass-bg);
  }

  .tab.active {
    background: var(--glass-bg);
    border-color: var(--glass-border);
    color: var(--text-primary);
  }

  .content-card {
    padding: 32px;
  }

  .content-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24px;
  }

  .loading-state, .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 64px;
    gap: 16px;
    color: var(--text-muted);
  }

  .empty-icon {
    font-size: 48px;
    opacity: 0.5;
  }

  .stock-control {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .stock-value {
    min-width: 40px;
    text-align: center;
    font-weight: 600;
  }

  .stock-value.low {
    color: var(--warning);
  }

  .action-buttons {
    display: flex;
    gap: 8px;
  }

  .categories-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
  }

  .category-card {
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    cursor: default;
  }

  .category-info h4 {
    margin-bottom: 4px;
  }

  .category-info p {
    font-size: 14px;
    margin-bottom: 8px;
  }

  .category-actions {
    display: flex;
    gap: 8px;
  }

  code {
    background: rgba(255, 255, 255, 0.1);
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 13px;
  }
</style>

