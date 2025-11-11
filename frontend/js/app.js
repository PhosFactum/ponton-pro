// frontend/js/app.js
const API_BASE = 'http://localhost:8080/api';

document.addEventListener('DOMContentLoaded', () => {
  initUI();
  fetchProducts();
});

function initUI() {
  // Модал и форма
  const modal = document.getElementById('orderModal');
  const orderForm = document.getElementById('orderForm');
  const closeBtn = document.querySelector('.close');

  // Закрыть по кресту
  closeBtn.addEventListener('click', closeOrderModal);

  // Закрыть по клику вне контента
  modal.addEventListener('click', (e) => {
    if (e.target === modal) closeOrderModal();
  });

  // Esc — закрыть
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') closeOrderModal();
  });

  // Обработка отправки формы
  orderForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    await submitOrder();
  });
}

async function fetchProducts() {
  const list = document.getElementById('productsList');
  list.innerHTML = '<div class="loading">Загрузка продуктов...</div>';

  try {
    const res = await fetch(`${API_BASE}/products`, {
      method: 'GET',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
    });

    if (!res.ok) throw new Error(`HTTP ${res.status} ${res.statusText}`);
    const payload = await res.json();

    // Поддерживаем несколько форматов ответа (на всякий случай)
    const products = payload.data ?? payload.products ?? payload;
    if (!Array.isArray(products) || products.length === 0) {
      list.innerHTML = '<div class="loading">Продуктов не найдено.</div>';
      return;
    }

    displayProducts(products);
  } catch (err) {
    console.error('Ошибка загрузки продуктов:', err);
    document.getElementById('productsList').innerHTML =
      `<div class="loading">Ошибка загрузки: ${escapeHtml(err.message)}</div>`;
  }
}

function displayProducts(products) {
  const html = products.map(p => productCardHtml(p)).join('');
  document.getElementById('productsList').innerHTML = html;

  // Привязываем события "Заказать" (вешаем делегированно)
  document.getElementById('productsList').addEventListener('click', (e) => {
    const btn = e.target.closest('[data-order-id]');
    if (!btn) return;
    openOrderModal(btn.getAttribute('data-order-id'), btn.getAttribute('data-order-title'));
  }, { once: false });
}

function productCardHtml(product) {
  // Поля, которые могут быть у продукта: id, title, name, description, price, image
  const id = product.id ?? product.ID ?? product.Id ?? '';
  const title = escapeHtml(product.title ?? product.name ?? 'Без названия');
  const description = escapeHtml(product.description ?? '');
  const price = product.price !== undefined ? `Цена: ${escapeHtml(String(product.price))}` : '';

  return `
    <div class="product-card">
      <h3>${title}</h3>
      <p>${description}</p>
      <p style="font-weight:600">${price}</p>
      <button class="order-btn" data-order-id="${id}" data-order-title="${title}">Оставить заявку</button>
    </div>
  `;
}

function openOrderModal(productId, productTitle = '') {
  const modal = document.getElementById('orderModal');
  modal.style.display = 'block';
  document.getElementById('productId').value = productId || '';
  // Подставим в заголовок модалки (если есть)
  const heading = modal.querySelector('h3');
  if (heading) heading.textContent = productTitle ? `Оставить заявку — ${productTitle}` : 'Оставить заявку';
  // Фокус на имя
  setTimeout(() => {
    const nameInput = document.getElementById('name');
    if (nameInput) nameInput.focus();
  }, 50);
}

function closeOrderModal() {
  const modal = document.getElementById('orderModal');
  modal.style.display = 'none';
  // очистим форму
  document.getElementById('orderForm').reset();
  document.getElementById('productId').value = '';
}

async function submitOrder() {
  const productId = document.getElementById('productId').value;
  const name = document.getElementById('name').value.trim();
  const phone = document.getElementById('phone').value.trim();
  const email = document.getElementById('email').value.trim();
  const description = document.getElementById('description').value.trim();

  if (!name || !phone) {
    showNotification('Имя и телефон — обязательны', 'error');
    return;
  }

  const body = {
    product_id: productId,
    name,
    phone,
    email: email || undefined,
    description: description || undefined
  };

  try {
    const res = await fetch(`${API_BASE}/requests`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    });

    if (!res.ok) {
      // Попробуем получить json с ошибкой
      let errMsg = `HTTP ${res.status}`;
      try {
        const errJson = await res.json();
        if (errJson && errJson.error) errMsg = errJson.error;
        else if (errJson && typeof errJson === 'string') errMsg = errJson;
      } catch (_) { /* ignore */ }
      throw new Error(errMsg);
    }

    let data;
    try { data = await res.json(); } catch (_) { data = null; }
    showNotification('Заявка успешно отправлена!', 'success');
    closeOrderModal();
    console.log('Ответ сервера на создание заявки:', data);
  } catch (err) {
    console.error('Ошибка отправки заявки:', err);
    showNotification(`Ошибка отправки: ${escapeHtml(err.message)}`, 'error');
  }
}

function showNotification(text, type = 'success', timeout = 4000) {
  const n = document.getElementById('notification');
  n.textContent = text;
  n.className = `notification ${type === 'error' ? 'error' : 'success'}`;
  n.style.display = 'block';
  if (n._hideTimeout) clearTimeout(n._hideTimeout);
  n._hideTimeout = setTimeout(() => {
    n.style.display = 'none';
  }, timeout);
}

// Простая защита от XSS при вставке текста
function escapeHtml(str) {
  if (str === null || str === undefined) return '';
  return String(str)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}
