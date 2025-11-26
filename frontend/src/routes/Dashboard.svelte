<script>
  import { onMount, onDestroy } from 'svelte';
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';
  import { flip } from 'svelte/animate';
  import { fade, slide, fly } from 'svelte/transition';
  import { categories as categoriesAPI, products as productsAPI } from '../lib/api.js';
  import { notifications } from '../lib/stores/notifications.js';
  import { websocketStore } from '../lib/stores/websocket.js';
  import { isAdmin } from '../lib/stores/auth.js';
  import Navbar from '../lib/components/Navbar.svelte';
  import Modal from '../lib/components/Modal.svelte';

  // State
  let categoriesData = { data: [], total_items: 0 };
  let productsData = { data: [], total_items: 0, page: 1, page_size: 20, total_pages: 1 };
  let isLoading = true;
  let activeTab = 'products';

  // Track changed rows for flash effect with direction
  let changedProducts = new Map();

  // Pagination & Search (server-side)
  let currentPage = 1;
  let pageSize = 20;
  let searchQuery = '';
  let selectedCategory = 'all';
  let searchTimeout = null;

  // Sorting (client-side for current page)
  let sortColumn = 'name';
  let sortDirection = 'asc';

  // Modal state
  let showCategoryModal = false;
  let showProductModal = false;
  let showLowStockModal = false;
  let showImportModal = false;
  let editingCategory = null;
  let editingProduct = null;

  // Form data
  let categoryForm = { name: '', description: '' };
  let productForm = { name: '', description: '', sku: '', quantity: 0, price: 0, category_id: 0 };

  // Import JSON
  let importData = '';
  let importErrors = [];
  let isImporting = false;
  let fileInput;

  // Tweened values for animated counters
  const tweenedTotalProducts = tweened(0, { duration: 800, easing: cubicOut });
  const tweenedTotalCategories = tweened(0, { duration: 800, easing: cubicOut });
  const tweenedTotalStock = tweened(0, { duration: 800, easing: cubicOut });
  const tweenedInventoryValue = tweened(0, { duration: 800, easing: cubicOut });

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

  function flashProduct(id, direction = 'neutral') {
    changedProducts.set(id, direction);
    changedProducts = changedProducts;
    setTimeout(() => {
      changedProducts.delete(id);
      changedProducts = changedProducts;
    }, 2000);
  }

  function handleProductCreated(product) {
    productsData = {
      ...productsData,
      data: [product, ...productsData.data],
      total_items: productsData.total_items + 1,
    };
    flashProduct(product.id, 'up');
    updateTweenedValues();
  }

  function handleProductUpdated(product) {
    const oldProduct = productsData.data.find(p => p.id === product.id);
    let direction = 'neutral';
    if (oldProduct) {
      if (product.quantity > oldProduct.quantity || product.price > oldProduct.price) {
        direction = 'up';
      } else if (product.quantity < oldProduct.quantity || product.price < oldProduct.price) {
        direction = 'down';
      }
    }
    productsData = {
      ...productsData,
      data: productsData.data.map(p => p.id === product.id ? product : p),
    };
    flashProduct(product.id, direction);
    updateTweenedValues();
  }

  function handleProductDeleted(payload) {
    productsData = {
      ...productsData,
      data: productsData.data.filter(p => p.id !== payload.id),
      total_items: productsData.total_items - 1,
    };
    updateTweenedValues();
  }

  function handleStockUpdated(product) {
    handleProductUpdated(product);
  }

  function updateTweenedValues() {
    tweenedTotalProducts.set(productsData.total_items);
    tweenedTotalCategories.set(categoriesData.total_items);
    const stock = productsData.data.reduce((sum, p) => sum + p.quantity, 0);
    const value = productsData.data.reduce((sum, p) => sum + (p.price * p.quantity), 0);
    tweenedTotalStock.set(stock);
    tweenedInventoryValue.set(value);
  }

  async function loadData(resetPage = true) {
    isLoading = true;
    try {
      if (resetPage) currentPage = 1;

      const categoryId = selectedCategory !== 'all' ? parseInt(selectedCategory) : null;

      const [cats, prods] = await Promise.all([
        categoriesAPI.list(1, 100),
        productsAPI.list(currentPage, pageSize, categoryId, searchQuery),
      ]);
      categoriesData = cats;
      productsData = prods;
      tweenedTotalProducts.set(prods.total_items);
      tweenedTotalCategories.set(cats.total_items);
      const stock = prods.data.reduce((sum, p) => sum + p.quantity, 0);
      const value = prods.data.reduce((sum, p) => sum + (p.price * p.quantity), 0);
      tweenedTotalStock.set(stock);
      tweenedInventoryValue.set(value);
    } catch (err) {
      notifications.error('Failed to load data');
    } finally {
      isLoading = false;
    }
  }

  // Load products with current filters (for pagination)
  async function loadProducts(resetPage = true) {
    isLoading = true;
    try {
      if (resetPage) currentPage = 1;

      const categoryId = selectedCategory !== 'all' ? parseInt(selectedCategory) : null;
      const prods = await productsAPI.list(currentPage, pageSize, categoryId, searchQuery);
      productsData = prods;
      tweenedTotalProducts.set(prods.total_items);
      const stock = prods.data.reduce((sum, p) => sum + p.quantity, 0);
      const value = prods.data.reduce((sum, p) => sum + (p.price * p.quantity), 0);
      tweenedTotalStock.set(stock);
      tweenedInventoryValue.set(value);
    } catch (err) {
      notifications.error('Failed to load products');
    } finally {
      isLoading = false;
    }
  }

  // Debounced search
  function handleSearch(value) {
    searchQuery = value;
    if (searchTimeout) clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      loadProducts(true);
    }, 300);
  }

  // Category filter change
  function handleCategoryChange(value) {
    selectedCategory = value;
    loadProducts(true);
  }

  // Pagination
  function goToPage(page) {
    if (page < 1 || page > productsData.total_pages) return;
    currentPage = page;
    loadProducts(false);
  }

  function nextPage() {
    goToPage(currentPage + 1);
  }

  function prevPage() {
    goToPage(currentPage - 1);
  }

  // Calculate which page numbers to show
  function getPageNumber(index) {
    const totalPages = productsData.total_pages;
    if (totalPages <= 5) return index + 1;

    // Show pages around current page
    let start = Math.max(1, currentPage - 2);
    let end = Math.min(totalPages, start + 4);

    if (end - start < 4) {
      start = Math.max(1, end - 4);
    }

    const page = start + index;
    return page <= end ? page : 0;
  }

  function handleSort(column) {
    if (sortColumn === column) {
      sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
      sortColumn = column;
      sortDirection = 'asc';
    }
  }

  // Sort products client-side (filtering is now server-side)
  $: sortedProducts = [...productsData.data].sort((a, b) => {
    let comparison = 0;
    switch (sortColumn) {
      case 'name': comparison = a.name.localeCompare(b.name); break;
      case 'sku': comparison = a.sku.localeCompare(b.sku); break;
      case 'stock': comparison = a.quantity - b.quantity; break;
      case 'price': comparison = a.price - b.price; break;
    }
    return sortDirection === 'asc' ? comparison : -comparison;
  });

  function exportToCSV() {
    const headers = ['Name', 'SKU', 'Category', 'Quantity', 'Price', 'Value'];
    const rows = sortedProducts.map(p => [
      `"${p.name}"`, p.sku, `"${getCategoryName(p.category_id)}"`,
      p.quantity, p.price.toFixed(2), (p.price * p.quantity).toFixed(2)
    ]);
    const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.setAttribute('href', url);
    link.setAttribute('download', `inventory_${new Date().toISOString().split('T')[0]}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    notifications.success('Exported to CSV');
  }

  function exportToJSON() {
    const data = {
      exported_at: new Date().toISOString(),
      categories: categoriesData.data.map(c => ({
        name: c.name,
        description: c.description
      })),
      products: productsData.data.map(p => ({
        name: p.name,
        description: p.description,
        sku: p.sku,
        quantity: p.quantity,
        price: p.price,
        category: getCategoryName(p.category_id)
      }))
    };
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.setAttribute('href', url);
    link.setAttribute('download', `inventory_${new Date().toISOString().split('T')[0]}.json`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    notifications.success('Exported to JSON');
  }

  function openImportModal() {
    importData = '';
    importErrors = [];
    showImportModal = true;
  }

  function handleFileSelect(event) {
    const file = event.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      importData = e.target.result;
    };
    reader.readAsText(file);
  }

  async function handleImport() {
    importErrors = [];
    isImporting = true;

    try {
      const data = JSON.parse(importData);

      // Validate structure
      if (!data.products && !data.categories) {
        importErrors = ['Invalid JSON format. Expected "products" or "categories" array.'];
        isImporting = false;
        return;
      }

      let imported = { categories: 0, products: 0 };

      // Import categories first
      if (data.categories && Array.isArray(data.categories)) {
        for (const cat of data.categories) {
          if (!cat.name) {
            importErrors.push(`Category missing name: ${JSON.stringify(cat)}`);
            continue;
          }
          try {
            // Check if category exists
            const existing = categoriesData.data.find(c => c.name.toLowerCase() === cat.name.toLowerCase());
            if (!existing) {
              await categoriesAPI.create({ name: cat.name, description: cat.description || '' });
              imported.categories++;
            }
          } catch (err) {
            importErrors.push(`Failed to import category "${cat.name}": ${err.message}`);
          }
        }
        // Reload categories
        const cats = await categoriesAPI.list(1, 100);
        categoriesData = cats;
      }

      // Import products
      if (data.products && Array.isArray(data.products)) {
        for (const prod of data.products) {
          if (!prod.name || !prod.sku) {
            importErrors.push(`Product missing name or sku: ${JSON.stringify(prod)}`);
            continue;
          }
          try {
            // Find category by name
            let categoryId = categoriesData.data[0]?.id;
            if (prod.category) {
              const cat = categoriesData.data.find(c => c.name.toLowerCase() === prod.category.toLowerCase());
              if (cat) categoryId = cat.id;
            } else if (prod.category_id) {
              categoryId = prod.category_id;
            }

            // Check if product with same SKU exists
            const existing = productsData.data.find(p => p.sku.toLowerCase() === prod.sku.toLowerCase());
            if (existing) {
              // Update existing
              await productsAPI.update(existing.id, {
                name: prod.name,
                description: prod.description || '',
                sku: prod.sku,
                quantity: prod.quantity || 0,
                price: prod.price || 0,
                category_id: categoryId
              });
            } else {
              // Create new
              await productsAPI.create({
                name: prod.name,
                description: prod.description || '',
                sku: prod.sku,
                quantity: prod.quantity || 0,
                price: prod.price || 0,
                category_id: categoryId
              });
            }
            imported.products++;
          } catch (err) {
            importErrors.push(`Failed to import product "${prod.name}": ${err.message}`);
          }
        }
      }

      // Reload data
      await loadData();

      if (importErrors.length === 0) {
        notifications.success(`Imported ${imported.categories} categories and ${imported.products} products`);
        showImportModal = false;
      } else {
        notifications.warning(`Import completed with ${importErrors.length} errors`);
      }
    } catch (err) {
      importErrors = [`Invalid JSON: ${err.message}`];
    } finally {
      isImporting = false;
    }
  }

  function openCategoryModal(category = null) {
    editingCategory = category;
    categoryForm = category ? { name: category.name, description: category.description } : { name: '', description: '' };
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

  function openProductModal(product = null) {
    editingProduct = product;
    productForm = product ? { ...product } : { name: '', description: '', sku: '', quantity: 0, price: 0, category_id: categoriesData.data[0]?.id || 0 };
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

  function getCategoryName(id) {
    const cat = categoriesData.data.find(c => c.id === id);
    return cat?.name || 'Unknown';
  }

  function getStockPercent(quantity) {
    return Math.min(quantity, 100);
  }

  function getStockStatus(quantity) {
    if (quantity <= 10) return 'critical';
    if (quantity <= 30) return 'warning';
    return 'good';
  }

  $: totalStock = productsData.data.reduce((sum, p) => sum + p.quantity, 0);
  $: inventoryValue = productsData.data.reduce((sum, p) => sum + (p.price * p.quantity), 0);
  $: lowStockProducts = productsData.data.filter(p => p.quantity < 10);
  $: lowStockCount = lowStockProducts.length;
  $: criticalStockCount = productsData.data.filter(p => p.quantity <= 3).length;
  $: avgPrice = productsData.data.length > 0 ? productsData.data.reduce((sum, p) => sum + p.price, 0) / productsData.data.length : 0;

  function formatNumber(n) {
    return Math.round(n).toLocaleString();
  }
</script>

<Navbar />

<main class="dashboard">
  <!-- KPIs HORIZONTAL - TOP BAR (FIXED) -->
  <header class="kpi-bar">
    <div class="kpi-card">
      <div class="kpi-icon blue">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
        </svg>
      </div>
      <div class="kpi-info">
        <span class="kpi-value">{formatNumber($tweenedTotalProducts)}</span>
        <span class="kpi-label">Products</span>
      </div>
      <span class="kpi-trend up">+3</span>
    </div>

    <div class="kpi-card">
      <div class="kpi-icon purple">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
        </svg>
      </div>
      <div class="kpi-info">
        <span class="kpi-value">{formatNumber($tweenedTotalCategories)}</span>
        <span class="kpi-label">Categories</span>
      </div>
    </div>

    <div class="kpi-card">
      <div class="kpi-icon" class:orange={lowStockCount > 0} class:teal={lowStockCount === 0}>
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
        </svg>
      </div>
      <div class="kpi-info">
        <span class="kpi-value">{formatNumber($tweenedTotalStock)}</span>
        <span class="kpi-label">Total Stock</span>
      </div>
      {#if lowStockCount > 0}
        <button class="kpi-alert" on:click={() => showLowStockModal = true}>
          {lowStockCount} low
        </button>
      {/if}
    </div>

    <div class="kpi-card highlight">
      <div class="kpi-icon green">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
      </div>
      <div class="kpi-info">
        <span class="kpi-value">${($tweenedInventoryValue / 1000).toFixed(1)}k</span>
        <span class="kpi-label">Inventory Value</span>
      </div>
      <span class="kpi-trend up">+12%</span>
    </div>
  </header>

  <!-- MAIN CONTENT -->
  <div class="content-wrapper">
    <!-- TABLE CARD -->
    <div class="table-card">
      <!-- Table Header (FIXED) -->
      <div class="table-header">
        <div class="header-left">
          <div class="tabs">
            <button class="tab {activeTab === 'products' ? 'active' : ''}" on:click={() => activeTab = 'products'}>
              Products
              {#if lowStockCount > 0}
                <span class="badge-alert">{lowStockCount}</span>
              {/if}
            </button>
            <button class="tab {activeTab === 'categories' ? 'active' : ''}" on:click={() => activeTab = 'categories'}>
              Categories
            </button>
          </div>
        </div>
        <div class="header-right">
          {#if activeTab === 'products'}
            <div class="search-box">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
              </svg>
              <input type="text" placeholder="Search products..." value={searchQuery} on:input={(e) => handleSearch(e.target.value)} />
            </div>
            <select class="filter-select" value={selectedCategory} on:change={(e) => handleCategoryChange(e.target.value)}>
              <option value="all">All Categories</option>
              {#each categoriesData.data as cat}
                <option value={cat.id}>{cat.name}</option>
              {/each}
            </select>
            <button class="btn-icon" on:click={openImportModal} title="Import JSON">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
            </button>
            <button class="btn-icon" on:click={exportToCSV} title="Export CSV">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
            </button>
            <button class="btn-icon" on:click={exportToJSON} title="Export JSON">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
              </svg>
            </button>
          {/if}
          {#if $isAdmin}
            <button class="btn-primary" on:click={() => activeTab === 'products' ? openProductModal() : openCategoryModal()}>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              Add {activeTab === 'products' ? 'Product' : 'Category'}
            </button>
          {/if}
        </div>
      </div>

      <!-- Table Content (SCROLLABLE) -->
      <div class="table-body">
        {#if isLoading}
          <div class="skeleton-table">
            {#each Array(6) as _, i}
              <div class="skeleton-row" style="animation-delay: {i * 0.1}s">
                <div class="skeleton-cell wide"></div>
                <div class="skeleton-cell"></div>
                <div class="skeleton-cell"></div>
                <div class="skeleton-cell"></div>
                <div class="skeleton-cell narrow"></div>
              </div>
            {/each}
          </div>
        {:else if activeTab === 'products'}
          {#if sortedProducts.length === 0}
            <div class="empty-state">
              <div class="empty-icon">
                <svg width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                  <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
                </svg>
              </div>
              <h3>{searchQuery || selectedCategory !== 'all' ? 'No matching products' : 'No products yet'}</h3>
              <p>{searchQuery || selectedCategory !== 'all' ? 'Try adjusting your search or filters' : 'Get started by adding your first product'}</p>
              {#if $isAdmin && !searchQuery && selectedCategory === 'all'}
                <button class="btn-primary" on:click={() => openProductModal()}>
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
                  </svg>
                  Add Product
                </button>
              {/if}
            </div>
          {:else}
            <div class="table-scroll">
              <table>
                <thead>
                  <tr>
                    <th class="sortable" class:active={sortColumn === 'name'} on:click={() => handleSort('name')}>
                      Product
                      {#if sortColumn === 'name'}<span class="sort-arrow">{sortDirection === 'asc' ? '↑' : '↓'}</span>{/if}
                    </th>
                    <th class="sortable" class:active={sortColumn === 'sku'} on:click={() => handleSort('sku')}>
                      SKU
                      {#if sortColumn === 'sku'}<span class="sort-arrow">{sortDirection === 'asc' ? '↑' : '↓'}</span>{/if}
                    </th>
                    <th class="sortable" class:active={sortColumn === 'stock'} on:click={() => handleSort('stock')}>
                      Stock
                      {#if sortColumn === 'stock'}<span class="sort-arrow">{sortDirection === 'asc' ? '↑' : '↓'}</span>{/if}
                    </th>
                    <th class="sortable text-right" class:active={sortColumn === 'price'} on:click={() => handleSort('price')}>
                      Price
                      {#if sortColumn === 'price'}<span class="sort-arrow">{sortDirection === 'asc' ? '↑' : '↓'}</span>{/if}
                    </th>
                    {#if $isAdmin}<th class="text-right">Actions</th>{/if}
                  </tr>
                </thead>
                <tbody>
                  {#each sortedProducts as product (product.id)}
                    {@const flashType = changedProducts.get(product.id)}
                    {@const stockStatus = getStockStatus(product.quantity)}
                    <tr
                      class:flash-up={flashType === 'up'}
                      class:flash-down={flashType === 'down'}
                      class:flash-neutral={flashType === 'neutral'}
                      animate:flip={{ duration: 300 }}
                      in:fly={{ y: -20, duration: 300 }}
                      out:slide={{ duration: 200 }}
                    >
                      <td>
                        <div class="product-info">
                          <span class="product-name">{product.name}</span>
                          <span class="product-category">{getCategoryName(product.category_id)}</span>
                        </div>
                      </td>
                      <td><code class="sku">{product.sku}</code></td>
                      <td>
                        <div class="stock-cell">
                          <span class="stock-value" class:critical={stockStatus === 'critical'} class:warning={stockStatus === 'warning'}>
                            {product.quantity}
                          </span>
                          <div class="stock-bar">
                            <div class="stock-fill {stockStatus}" style="width: {getStockPercent(product.quantity)}%"></div>
                          </div>
                        </div>
                      </td>
                      <td class="text-right">
                        <span class="price">${product.price.toFixed(2)}</span>
                      </td>
                      {#if $isAdmin}
                        <td class="text-right">
                          <div class="actions">
                            <button class="action-btn" on:click={() => openProductModal(product)} title="Edit">
                              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                              </svg>
                            </button>
                            <button class="action-btn danger" on:click={() => deleteProduct(product.id)} title="Delete">
                              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="3 6 5 6 21 6"/>
                                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                              </svg>
                            </button>
                          </div>
                        </td>
                      {/if}
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        {:else}
          <!-- Categories Tab -->
          {#if categoriesData.data.length === 0}
            <div class="empty-state">
              <div class="empty-icon">
                <svg width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                  <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
                </svg>
              </div>
              <h3>No categories yet</h3>
              <p>Create categories to organize your products</p>
              {#if $isAdmin}
                <button class="btn-primary" on:click={() => openCategoryModal()}>
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
                  </svg>
                  Add Category
                </button>
              {/if}
            </div>
          {:else}
            <div class="categories-grid">
              {#each categoriesData.data as category (category.id)}
                {@const productCount = productsData.data.filter(p => p.category_id === category.id).length}
                <div class="category-card" animate:flip={{ duration: 300 }} transition:fade={{ duration: 200 }}>
                  <div class="category-icon">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                      <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
                    </svg>
                  </div>
                  <div class="category-info">
                    <h4>{category.name}</h4>
                    <p>{category.description || 'No description'}</p>
                    <span class="product-count">{productCount} product{productCount !== 1 ? 's' : ''}</span>
                  </div>
                  {#if $isAdmin}
                    <div class="category-actions">
                      <button class="action-btn" on:click={() => openCategoryModal(category)} title="Edit">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                      </button>
                      <button class="action-btn danger" on:click={() => deleteCategory(category.id)} title="Delete">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <polyline points="3 6 5 6 21 6"/>
                          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                        </svg>
                      </button>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {/if}
        {/if}
      </div>

      <!-- Table Footer with Pagination (FIXED) -->
      {#if !isLoading && activeTab === 'products'}
        <div class="table-footer">
          <span class="footer-info">
            Showing {((currentPage - 1) * pageSize) + 1}–{Math.min(currentPage * pageSize, productsData.total_items)} of {productsData.total_items} products
          </span>
          {#if productsData.total_pages > 1}
            <div class="pagination">
              <button class="page-btn" on:click={prevPage} disabled={currentPage === 1} title="Previous">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="15 18 9 12 15 6"/>
                </svg>
              </button>

              {#each Array(Math.min(5, productsData.total_pages)) as _, i}
                {@const pageNum = getPageNumber(i)}
                {#if pageNum > 0}
                  <button
                    class="page-btn"
                    class:active={currentPage === pageNum}
                    on:click={() => goToPage(pageNum)}
                  >
                    {pageNum}
                  </button>
                {/if}
              {/each}

              {#if productsData.total_pages > 5 && currentPage < productsData.total_pages - 2}
                <span class="page-ellipsis">...</span>
                <button class="page-btn" on:click={() => goToPage(productsData.total_pages)}>
                  {productsData.total_pages}
                </button>
              {/if}

              <button class="page-btn" on:click={nextPage} disabled={currentPage === productsData.total_pages} title="Next">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="9 18 15 12 9 6"/>
                </svg>
              </button>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
</main>

<!-- Modals -->
<Modal title="Low Stock Alert" show={showLowStockModal} on:close={() => showLowStockModal = false}>
  <div class="low-stock-list">
    {#each lowStockProducts.sort((a, b) => a.quantity - b.quantity) as product (product.id)}
      <div class="low-stock-item" class:critical={product.quantity <= 3}>
        <div class="low-stock-info">
          <strong>{product.name}</strong>
          <code class="sku">{product.sku}</code>
        </div>
        <div class="low-stock-qty">
          <span class="stock-value" class:critical={product.quantity <= 3}>{product.quantity}</span>
          <span class="text-muted">units</span>
        </div>
      </div>
    {/each}
  </div>
</Modal>

<Modal title="Import Data (JSON)" show={showImportModal} on:close={() => showImportModal = false}>
  <div class="import-modal">
    <p class="import-help">
      Upload a JSON file or paste JSON data below. Expected format:
    </p>
    <pre class="import-example">{`{
  "categories": [
    { "name": "Electronics", "description": "..." }
  ],
  "products": [
    { "name": "Phone", "sku": "PHN-001", "quantity": 10, "price": 599.99, "category": "Electronics" }
  ]
}`}</pre>

    <div class="import-file">
      <input type="file" accept=".json" bind:this={fileInput} on:change={handleFileSelect} style="display: none" />
      <button class="btn-secondary" on:click={() => fileInput.click()}>
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="17 8 12 3 7 8"/>
          <line x1="12" y1="3" x2="12" y2="15"/>
        </svg>
        Choose File
      </button>
    </div>

    <div class="form-group">
      <label for="import-data">Or paste JSON data:</label>
      <textarea id="import-data" bind:value={importData} rows="8" placeholder={'{ "products": [...] }'}></textarea>
    </div>

    {#if importErrors.length > 0}
      <div class="import-errors">
        <strong>Errors:</strong>
        <ul>
          {#each importErrors as error}
            <li>{error}</li>
          {/each}
        </ul>
      </div>
    {/if}

    <div class="form-actions">
      <button type="button" class="btn-secondary" on:click={() => showImportModal = false}>Cancel</button>
      <button class="btn-primary" on:click={handleImport} disabled={!importData || isImporting}>
        {#if isImporting}
          <span class="spinner"></span>
          Importing...
        {:else}
          Import
        {/if}
      </button>
    </div>
  </div>
</Modal>

<Modal title={editingCategory ? 'Edit Category' : 'New Category'} show={showCategoryModal} on:close={() => showCategoryModal = false}>
  <form on:submit|preventDefault={saveCategory}>
    <div class="form-group">
      <label for="cat-name">Name</label>
      <input type="text" id="cat-name" bind:value={categoryForm.name} placeholder="Category name" required />
    </div>
    <div class="form-group">
      <label for="cat-desc">Description</label>
      <input type="text" id="cat-desc" bind:value={categoryForm.description} placeholder="Optional description" />
    </div>
    <div class="form-actions">
      <button type="button" class="btn-secondary" on:click={() => showCategoryModal = false}>Cancel</button>
      <button type="submit" class="btn-primary">{editingCategory ? 'Update' : 'Create'}</button>
    </div>
  </form>
</Modal>

<Modal title={editingProduct ? 'Edit Product' : 'New Product'} show={showProductModal} on:close={() => showProductModal = false}>
  <form on:submit|preventDefault={saveProduct}>
    <div class="form-group">
      <label for="prod-name">Name</label>
      <input type="text" id="prod-name" bind:value={productForm.name} placeholder="Product name" required />
    </div>
    <div class="form-group">
      <label for="prod-sku">SKU</label>
      <input type="text" id="prod-sku" bind:value={productForm.sku} placeholder="SKU-001" required />
    </div>
    <div class="form-group">
      <label for="prod-desc">Description</label>
      <input type="text" id="prod-desc" bind:value={productForm.description} placeholder="Optional description" />
    </div>
    <div class="form-group">
      <label for="prod-cat">Category</label>
      <select id="prod-cat" bind:value={productForm.category_id}>
        {#each categoriesData.data as cat}<option value={cat.id}>{cat.name}</option>{/each}
      </select>
    </div>
    <div class="form-row">
      <div class="form-group">
        <label for="prod-qty">Quantity</label>
        <input type="number" id="prod-qty" bind:value={productForm.quantity} min="0" required />
      </div>
      <div class="form-group">
        <label for="prod-price">Price ($)</label>
        <input type="number" id="prod-price" bind:value={productForm.price} min="0.01" step="0.01" required />
      </div>
    </div>
    <div class="form-actions">
      <button type="button" class="btn-secondary" on:click={() => showProductModal = false}>Cancel</button>
      <button type="submit" class="btn-primary">{editingProduct ? 'Update' : 'Create'}</button>
    </div>
  </form>
</Modal>

<style>
  /* LAYOUT - Fixed height dashboard */
  .dashboard {
    height: calc(100vh - 64px);
    background: #F8FAFC;
    padding: 24px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  /* KPI BAR - HORIZONTAL TOP (FIXED) */
  .kpi-bar {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
    margin-bottom: 24px;
    flex-shrink: 0;
  }

  .kpi-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px 24px;
    background: white;
    border-radius: 16px;
    border: 1px solid #E2E8F0;
    box-shadow: 0 1px 3px rgba(0,0,0,0.04);
    transition: all 0.2s ease;
  }

  .kpi-card:hover {
    box-shadow: 0 4px 12px rgba(0,0,0,0.08);
    transform: translateY(-2px);
  }

  .kpi-card.highlight {
    background: linear-gradient(135deg, #ECFDF5 0%, #D1FAE5 100%);
    border-color: #A7F3D0;
  }

  .kpi-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .kpi-icon.blue { background: #EFF6FF; color: #3B82F6; }
  .kpi-icon.purple { background: #F5F3FF; color: #8B5CF6; }
  .kpi-icon.teal { background: #F0FDFA; color: #14B8A6; }
  .kpi-icon.orange { background: #FFF7ED; color: #F97316; }
  .kpi-icon.green { background: #ECFDF5; color: #10B981; }

  .kpi-info {
    flex: 1;
    min-width: 0;
  }

  .kpi-value {
    display: block;
    font-size: 1.75rem;
    font-weight: 700;
    color: #0F172A;
    line-height: 1.2;
    font-variant-numeric: tabular-nums;
  }

  .kpi-label {
    font-size: 13px;
    color: #64748B;
    font-weight: 500;
  }

  .kpi-trend {
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 12px;
    font-weight: 600;
    flex-shrink: 0;
  }

  .kpi-trend.up {
    background: #ECFDF5;
    color: #059669;
  }

  .kpi-trend.down {
    background: #FEF2F2;
    color: #DC2626;
  }

  .kpi-alert {
    padding: 6px 12px;
    background: #FEF3C7;
    color: #D97706;
    border: none;
    border-radius: 6px;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .kpi-alert:hover {
    background: #FDE68A;
  }

  /* CONTENT WRAPPER - Takes remaining space */
  .content-wrapper {
    flex: 1;
    min-height: 0;
    display: flex;
    gap: 24px;
  }

  /* TABLE CARD - Fixed height with internal scroll */
  .table-card {
    flex: 1;
    background: white;
    border-radius: 16px;
    border: 1px solid #E2E8F0;
    box-shadow: 0 1px 3px rgba(0,0,0,0.04);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    min-height: 0;
  }

  .table-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding: 16px 24px;
    border-bottom: 1px solid #F1F5F9;
    flex-wrap: wrap;
    flex-shrink: 0;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  /* TABLE BODY - Scrollable area */
  .table-body {
    flex: 1;
    overflow-y: auto;
    min-height: 0;
  }

  .table-scroll {
    height: 100%;
  }

  /* TABLE FOOTER - Fixed at bottom with pagination */
  .table-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 24px;
    border-top: 1px solid #F1F5F9;
    background: #FAFBFC;
    font-size: 13px;
    color: #64748B;
    flex-shrink: 0;
    gap: 16px;
  }

  .footer-info {
    white-space: nowrap;
  }

  .pagination {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .page-btn {
    min-width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    border: 1px solid #E2E8F0;
    border-radius: 8px;
    color: #475569;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .page-btn:hover:not(:disabled) {
    background: #F1F5F9;
    border-color: #CBD5E1;
  }

  .page-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .page-btn.active {
    background: #10B981;
    border-color: #10B981;
    color: white;
  }

  .page-ellipsis {
    padding: 0 8px;
    color: #94A3B8;
  }

  /* TABS */
  .tabs {
    display: flex;
    gap: 4px;
    background: #F1F5F9;
    padding: 4px;
    border-radius: 10px;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 18px;
    background: transparent;
    border: none;
    color: #64748B;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.15s ease;
  }

  .tab:hover {
    color: #0F172A;
  }

  .tab.active {
    background: white;
    color: #0F172A;
    box-shadow: 0 1px 3px rgba(0,0,0,0.08);
  }

  .badge-alert {
    background: #F97316;
    color: white;
    font-size: 11px;
    font-weight: 600;
    padding: 2px 7px;
    border-radius: 100px;
  }

  /* SEARCH & FILTERS */
  .search-box {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 14px;
    background: #F8FAFC;
    border: 1px solid #E2E8F0;
    border-radius: 10px;
    width: 220px;
    transition: all 0.15s ease;
  }

  .search-box:focus-within {
    border-color: #10B981;
    background: white;
    box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
  }

  .search-box svg {
    color: #94A3B8;
    flex-shrink: 0;
  }

  .search-box input {
    flex: 1;
    border: none;
    background: transparent;
    font-size: 14px;
    color: #0F172A;
    outline: none;
  }

  .search-box input::placeholder {
    color: #94A3B8;
  }

  .filter-select {
    padding: 10px 36px 10px 14px;
    background: #F8FAFC;
    border: 1px solid #E2E8F0;
    border-radius: 10px;
    font-size: 14px;
    color: #0F172A;
    cursor: pointer;
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%2394A3B8' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 12px center;
  }

  .filter-select:focus {
    border-color: #10B981;
    outline: none;
    box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
  }

  /* BUTTONS */
  .btn-icon {
    width: 42px;
    height: 42px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #F8FAFC;
    border: 1px solid #E2E8F0;
    border-radius: 10px;
    color: #64748B;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .btn-icon:hover {
    background: #F1F5F9;
    color: #0F172A;
  }

  .btn-primary {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 18px;
    background: #10B981;
    color: white;
    border: none;
    border-radius: 10px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;
    box-shadow: 0 2px 8px rgba(16, 185, 129, 0.3);
  }

  .btn-primary:hover {
    background: #059669;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
  }

  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }

  .btn-secondary {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 18px;
    background: #F1F5F9;
    color: #475569;
    border: 1px solid #E2E8F0;
    border-radius: 10px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .btn-secondary:hover {
    background: #E2E8F0;
  }

  /* TABLE */
  table {
    width: 100%;
    border-collapse: collapse;
  }

  th {
    padding: 14px 24px;
    text-align: left;
    font-size: 12px;
    font-weight: 600;
    color: #64748B;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    background: #FAFBFC;
    border-bottom: 1px solid #F1F5F9;
    position: sticky;
    top: 0;
    z-index: 5;
  }

  th.sortable {
    cursor: pointer;
    user-select: none;
  }

  th.sortable:hover {
    color: #0F172A;
  }

  th.active {
    color: #10B981;
  }

  .sort-arrow {
    margin-left: 4px;
    font-size: 10px;
  }

  td {
    padding: 20px 24px;
    border-bottom: 1px solid #F1F5F9;
    vertical-align: middle;
  }

  tbody tr {
    transition: background 0.15s ease;
  }

  tbody tr:hover {
    background: #FAFBFC;
  }

  tbody tr:last-child td {
    border-bottom: none;
  }

  /* Product Info */
  .product-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .product-name {
    font-weight: 500;
    color: #0F172A;
    font-size: 15px;
  }

  .product-category {
    font-size: 12px;
    color: #94A3B8;
  }

  /* SKU */
  .sku {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
    color: #64748B;
    background: #F8FAFC;
    padding: 4px 8px;
    border-radius: 4px;
  }

  /* Stock Cell */
  .stock-cell {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 120px;
  }

  .stock-value {
    font-weight: 600;
    font-size: 15px;
    color: #0F172A;
    font-variant-numeric: tabular-nums;
  }

  .stock-value.warning {
    color: #D97706;
  }

  .stock-value.critical {
    color: #DC2626;
  }

  /* Stock Bar - Pills */
  .stock-bar {
    height: 8px;
    background: #F1F5F9;
    border-radius: 100px;
    overflow: hidden;
  }

  .stock-fill {
    height: 100%;
    border-radius: 100px;
    transition: width 0.5s ease;
  }

  .stock-fill.good {
    background: linear-gradient(90deg, #10B981, #34D399);
  }

  .stock-fill.warning {
    background: linear-gradient(90deg, #F59E0B, #FBBF24);
  }

  .stock-fill.critical {
    background: linear-gradient(90deg, #EF4444, #F87171);
  }

  /* Price */
  .price {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 15px;
    font-weight: 500;
    color: #0F172A;
  }

  /* Actions */
  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 6px;
  }

  .action-btn {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    color: #64748B;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.15s ease;
  }

  .action-btn:hover {
    background: #F1F5F9;
    color: #0F172A;
  }

  .action-btn.danger:hover {
    background: #FEF2F2;
    color: #DC2626;
  }

  /* Flash Animations */
  .flash-up {
    animation: flashGreen 2s ease-out;
  }

  .flash-down {
    animation: flashRed 2s ease-out;
  }

  .flash-neutral {
    animation: flashYellow 2s ease-out;
  }

  @keyframes flashGreen {
    0% { background: rgba(16, 185, 129, 0.15); }
    100% { background: transparent; }
  }

  @keyframes flashRed {
    0% { background: rgba(239, 68, 68, 0.15); }
    100% { background: transparent; }
  }

  @keyframes flashYellow {
    0% { background: rgba(245, 158, 11, 0.15); }
    100% { background: transparent; }
  }

  /* Empty State */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 64px 24px;
    text-align: center;
  }

  .empty-icon {
    color: #CBD5E1;
    margin-bottom: 16px;
  }

  .empty-state h3 {
    font-size: 18px;
    font-weight: 600;
    color: #0F172A;
    margin-bottom: 8px;
  }

  .empty-state p {
    color: #64748B;
    margin-bottom: 24px;
  }

  /* Skeleton */
  .skeleton-table {
    padding: 0 24px;
  }

  .skeleton-row {
    display: flex;
    gap: 24px;
    padding: 20px 0;
    border-bottom: 1px solid #F1F5F9;
    animation: shimmer 1.5s infinite;
  }

  .skeleton-cell {
    height: 20px;
    background: linear-gradient(90deg, #F1F5F9 25%, #E2E8F0 50%, #F1F5F9 75%);
    background-size: 200% 100%;
    border-radius: 4px;
    width: 15%;
  }

  .skeleton-cell.wide { width: 30%; }
  .skeleton-cell.narrow { width: 10%; }

  @keyframes shimmer {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }

  /* Categories Grid */
  .categories-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
    padding: 24px;
  }

  .category-card {
    display: flex;
    align-items: flex-start;
    gap: 16px;
    padding: 20px;
    background: #FAFBFC;
    border: 1px solid #F1F5F9;
    border-radius: 12px;
    transition: all 0.15s ease;
  }

  .category-card:hover {
    background: white;
    border-color: #E2E8F0;
    box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  }

  .category-icon {
    width: 44px;
    height: 44px;
    background: #EFF6FF;
    color: #3B82F6;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .category-info {
    flex: 1;
    min-width: 0;
  }

  .category-info h4 {
    font-size: 15px;
    font-weight: 600;
    color: #0F172A;
    margin-bottom: 4px;
  }

  .category-info p {
    font-size: 13px;
    color: #64748B;
    margin-bottom: 8px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .product-count {
    font-size: 12px;
    color: #94A3B8;
  }

  .category-actions {
    display: flex;
    gap: 4px;
  }

  /* Low Stock Modal */
  .low-stock-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    max-height: 400px;
    overflow-y: auto;
  }

  .low-stock-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 16px;
    background: #F8FAFC;
    border-radius: 10px;
    gap: 16px;
  }

  .low-stock-item.critical {
    background: #FEF2F2;
    border: 1px solid #FECACA;
  }

  .low-stock-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .low-stock-info strong {
    font-weight: 500;
    color: #0F172A;
  }

  .low-stock-qty {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  /* Import Modal */
  .import-modal {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .import-help {
    font-size: 14px;
    color: #64748B;
  }

  .import-example {
    background: #F8FAFC;
    border: 1px solid #E2E8F0;
    border-radius: 8px;
    padding: 12px;
    font-size: 12px;
    font-family: 'SF Mono', monospace;
    color: #475569;
    overflow-x: auto;
    white-space: pre;
  }

  .import-file {
    display: flex;
    justify-content: center;
    padding: 16px;
    border: 2px dashed #E2E8F0;
    border-radius: 10px;
  }

  .import-errors {
    background: #FEF2F2;
    border: 1px solid #FECACA;
    border-radius: 8px;
    padding: 12px;
    color: #DC2626;
    font-size: 13px;
  }

  .import-errors ul {
    margin: 8px 0 0 16px;
    padding: 0;
  }

  .import-errors li {
    margin-bottom: 4px;
  }

  .spinner {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255,255,255,0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .text-muted {
    color: #94A3B8;
    font-size: 13px;
  }

  .text-right {
    text-align: right;
  }

  /* Form */
  .form-group {
    margin-bottom: 16px;
  }

  .form-group label {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: #475569;
    margin-bottom: 6px;
  }

  .form-group input,
  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 12px 14px;
    background: #F8FAFC;
    border: 1px solid #E2E8F0;
    border-radius: 10px;
    font-size: 15px;
    font-family: inherit;
    color: #0F172A;
    transition: all 0.15s ease;
    resize: vertical;
  }

  .form-group input:focus,
  .form-group select:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: #10B981;
    background: white;
    box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
  }

  .form-group input::placeholder,
  .form-group textarea::placeholder {
    color: #94A3B8;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
    padding-top: 16px;
    border-top: 1px solid #F1F5F9;
  }

  /* RESPONSIVE */
  @media (max-width: 1200px) {
    .kpi-bar {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (max-width: 768px) {
    .dashboard {
      padding: 16px;
      height: auto;
      min-height: calc(100vh - 64px);
    }

    .kpi-bar {
      grid-template-columns: 1fr;
      gap: 12px;
    }

    .table-header {
      flex-direction: column;
      align-items: stretch;
      gap: 12px;
    }

    .header-left, .header-right {
      width: 100%;
    }

    .header-right {
      flex-wrap: wrap;
    }

    .search-box {
      width: 100%;
    }

    .filter-select {
      flex: 1;
    }

    td, th {
      padding: 14px 16px;
    }
  }
</style>
