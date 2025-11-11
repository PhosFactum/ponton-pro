const API_BASE = 'http://localhost:8080/api';

document.addEventListener('DOMContentLoaded', () => {
    initUI();
    fetchProducts();
});

function initUI() {
    const form = document.getElementById('orderForm');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            await submitOrder();
        });
    }

    const modal = document.getElementById('orderModal');
    if (modal) {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) closeOrderModal();
        });
    }

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') closeOrderModal();
    });
}

async function fetchProducts() {
    const list = document.getElementById('productsList');
    list.innerHTML = '<div class="loading">Загрузка продуктов...</div>';

    try {
        const res = await fetch(`${API_BASE}/products`, { method: 'GET', mode: 'cors' });
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const body = await res.json();

        const products = body.data ?? body.products ?? body;
        if (!Array.isArray(products) || products.length === 0) {
            list.innerHTML = '<div class="loading">Продуктов не найдено.</div>';
            return;
        }

        list.innerHTML = products.map(p => productHtml(p)).join('');
    } catch (err) {
        console.error('fetchProducts error:', err);
        list.innerHTML = `<div class="loading">Ошибка загрузки: ${escapeHtml(err.message)}</div>`;
    }
}

function productHtml(p) {
    const id = p.id ?? p.ID ?? p.Id ?? '';
    const title = escapeHtml(p.title ?? p.name ?? 'Без названия');
    const desc = escapeHtml(p.description ?? '');
    const imgIndex = id && !isNaN(id) ? id : 1;
    const img = `images/${imgIndex}.jpg`;

    return `
        <div class="product-card">
            <img src="${img}" alt="${title}" loading="lazy">
            <div class="product-info">
                <div>
                    <h3>${title}</h3>
                    <p>${desc}</p>
                </div>
                <button class="order-btn" type="button" onclick="openOrderModal('${id}', '${title.replace(/'/g, "\\'")}')">
                    Оставить заявку
                </button>
            </div>
        </div>
    `;
}

function openOrderModal(id, title) {
    document.getElementById('productId').value = id;
    const modal = document.getElementById('orderModal');
    if (!modal) return;

    modal.style.display = 'block';
    const h = modal.querySelector('h3');
    if (h) {
        h.textContent = title ? `Оставить заявку — ${title}` : 'Оставить заявку';
    }
    setTimeout(() => {
        const nameInput = document.getElementById('name');
        if (nameInput) nameInput.focus();
    }, 60);
}

function closeOrderModal() {
    const modal = document.getElementById('orderModal');
    if (modal) {
        modal.style.display = 'none';
        document.getElementById('orderForm')?.reset();
        document.getElementById('productId').value = '';
    }
}

async function submitOrder() {
    const product_id = Number(document.getElementById('productId').value) || undefined;
    const name = document.getElementById('name').value.trim();
    const phone = document.getElementById('phone').value.trim();
    const email = document.getElementById('email').value.trim();
    const description = document.getElementById('description').value.trim();

    if (!name || !phone) {
        showNotification('Имя и телефон обязательны', 'error');
        return;
    }

    const payload = { name, phone, product_id };
    if (email) payload.email = email;
    if (description) payload.description = description;

    try {
        const res = await fetch(`${API_BASE}/requests`, {
            method: 'POST',
            mode: 'cors',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        const data = await res.json();
        if (res.ok && (data.success || data.id)) {
            showNotification(data.message || 'Заявка успешно отправлена!', 'success');
            closeOrderModal();
        } else {
            showNotification(data.message || 'Ошибка при отправке заявки', 'error');
        }
    } catch (err) {
        showNotification(`Ошибка сети: ${err.message}`, 'error');
    }
}

function showNotification(msg, type = 'success', timeout = 4000) {
    const n = document.getElementById('notification');
    if (!n) return;

    n.textContent = msg;
    n.className = `notification ${type === 'error' ? 'error' : 'success'}`;
    n.style.display = 'block';

    if (n._t) clearTimeout(n._t);
    n._t = setTimeout(() => { n.style.display = 'none'; }, timeout);
}

function escapeHtml(str) {
    if (str == null) return '';
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '<')
        .replace(/>/g, '>')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}