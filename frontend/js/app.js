// ЗАМЕНИ ВЕСЬ файл app.js на этот упрощенный:
const API_BASE = 'http://localhost:8080/api';

document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 Page loaded, testing API...');
    testAPI();
});

async function testAPI() {
    try {
        console.log('🔍 Testing connection to:', API_BASE + '/products');
        
        const response = await fetch(API_BASE + '/products', {
            method: 'GET',
            mode: 'cors',
            headers: {
                'Content-Type': 'application/json',
            }
        });
        
        console.log('📡 Response status:', response.status);
        console.log('📡 Response headers:', response.headers);
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const data = await response.json();
        console.log('✅ SUCCESS! Data:', data);
        
        if (data && data.data) {
            displayProducts(data.data);
        } else {
            document.getElementById('productsList').innerHTML = 
                '<div class="error">Нет данных о продуктах</div>';
        }
        
    } catch (error) {
        console.error('❌ ERROR:', error);
        document.getElementById('productsList').innerHTML = 
            `<div class="error">Ошибка: ${error.message}</div>`;
    }
}

function displayProducts(products) {
    const html = products.map(product => `
        <div class="product-card">
            <h3>${product.title}</h3>
            <p>${product.description}</p>
            <button onclick="alert('Заказ продукта ${product.id}')">Заказать</button>
        </div>
    `).join('');
    
    document.getElementById('productsList').innerHTML = html;
}

// Временные функции для кнопок
function scrollToProducts() {
    const element = document.getElementById('products');
    if (element) {
        element.scrollIntoView({ behavior: 'smooth' });
        console.log('🔍 Scrolled to products');
    }
}

function openOrderModal(id) {
    alert('Открыть заказ продукта ' + id);
}

function closeOrderModal() {
    alert('Закрыть модальное окно');
}