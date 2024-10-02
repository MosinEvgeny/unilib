<template>
    <div class="container">
        <h1>UniLib - Каталог книг</h1>

        <div class="search-filters">
            <input type="text" v-model="searchQuery" placeholder="Поиск...">
            <button @click="searchBooks">Поиск</button>
            <!-- Дополнительные фильтры (по категории, году и т.д.) -->
        </div>

        <ul class="book-list">
            <li v-for="book in books" :key="book.book_id">
                <h3>{{ book.title }}</h3>
                <p>Автор: {{ book.author }}</p>
                <p v-if="book.availableCopies !== undefined">
                    Доступно: {{ book.availableCopies }} экз.
                </p>
                <p v-else-if="book.availableCopies === 0">Нет в наличии</p>
                <p v-else>Загрузка...</p>
            </li>
        </ul>

        <router-link to="/reader/history">История заказов</router-link>
        <button @click="logout">Выход</button>
    </div>
</template>

<script>
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

export default {
    setup() {
        const router = useRouter();
        const books = ref([]);
        const searchQuery = ref('');

        const getAvailableCopies = async (bookId) => {
            try {
                const response = await fetch(`http://localhost:8080/books/${bookId}/available`);
                const data = await response.json();
                return data.available_copies;
            } catch (error) {
                console.error(`Ошибка при получении доступных экземпляров для книги ${bookId}:`, error);
                return undefined; //  Возвращаем undefined в случае ошибки
            }
        };

        const logout = () => {
            //  Логика выхода (очистка данных авторизации, перенаправление на страницу входа)
            router.push('/');
        };

        onMounted(async () => {
            await searchBooks(); // Загружаем все книги при загрузке компонента
        });

        const searchBooks = async () => {
            try {
                const response = await fetch(`http://localhost:8080/books?search=${searchQuery.value}`);
                const data = await response.json();

                // Дожидаемся результата для каждой книги
                books.value = await Promise.all(data.map(async (book) => {
                    const availableCopies = await getAvailableCopies(book.book_id);
                    return {
                        ...book,
                        availableCopies,
                    };
                }));
            } catch (error) {
                console.error('Ошибка при поиске книг:', error);
            }
        };

        return {
            books,
            searchQuery,
            searchBooks,
            getAvailableCopies,
            logout
        };
    },
};
</script>

<style scoped></style>