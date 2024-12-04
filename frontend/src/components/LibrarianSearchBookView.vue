<template>
    <div class="container">
        <h1>Каталог книг</h1>

        <div class="search-filters">
            <input type="text" v-model="searchQuery" placeholder="Поиск...">
            <button @click="searchBooks">Поиск</button>
        </div>

        <div class="book-grid" v-if="books.length > 0">
            <div v-for="book in books" :key="book.book_id" class="book-item">
                <h3>{{ book.title }}</h3>
                <p>Автор: {{ book.author }}</p>
                <p v-if="book.availableCopies !== undefined">
                    Доступно: {{ book.availableCopies }} экз.
                </p>
                <p v-else-if="book.availableCopies === 0">Нет в наличии</p>
                <p v-else>Загрузка...</p>
            </div>
        </div>
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

<style scoped>
.search-filters {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
}

.book-list li {
    background-color: #f8f8f8;
}

.book-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1rem;
}

.book-item {
    /*  Добавляем  стили  для  .book-item  */
    border: 1px solid #eee;
    border-radius: 4px;
    padding: 1rem;
}
</style>