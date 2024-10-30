<template>
    <div class="container">
        <h1>История заказов</h1>

        <ul v-if="orders && orders.length > 0" class="order-list">
            <li v-for="order in orders" :key="order.issue_id">
                <h3>{{ order.title }}</h3>
                <p>Дата выдачи: {{ formatDate(order.issue_date) }}</p>
                <p>Срок возврата: {{ formatDate(order.due_date) }}</p>
                <p v-if="order.return_date">Дата возврата: {{ formatDate(order.return_date) }}</p>
                <p v-else-if="isOverdue(order.due_date)" style= "color:red">Просрочено</p>
                <p v-else style= "color:orange" >Не возвращена</p>
            </li>
        </ul>
        <p v-else>У вас нет заказов.</p>
    </div>
</template>

<script>
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

export default {
    setup() {
        const router = useRouter();
        const orders = ref([]);

        const fetchOrders = async () => {
            try {
                const readerId = localStorage.getItem('reader_id');
                const response = await fetch(`http://localhost:8080/orders/${readerId}`);
                orders.value = await response.json();
            } catch (error) {
                console.error('Ошибка при загрузке истории заказов:', error);
            }
        };

        const formatDate = (dateString) => {
            const date = new Date(dateString);
            return date.toLocaleDateString();
        };

        const isOverdue = (dueDate) => {
            const today = new Date();
            const due = new Date(dueDate);
            return due < today;
        };

        const goBack = () => {
            router.push('/reader');
        };

        onMounted(async () => {
            await fetchOrders();
        });

        return {
            orders,
            formatDate,
            isOverdue,
            goBack
        };
    },
};
</script>

<style scoped>
/* Стили для страницы истории заказов */
</style>