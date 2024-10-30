<template>
    <div class="container">
        <h1>Поиск книг</h1>

        <div class="search-filters">
            <input type="text" v-model="searchQuery" placeholder="Поиск по названию или автору...">
            <button @click="searchBooks">Поиск</button>
        </div>

        <ul class="book-list" v-if="books.length > 0">
            <li v-for="book in books" :key="book.book_id">
                <h3>{{ book.title }}</h3>
                <p>Автор: {{ book.author }}</p>
                <p>Доступно: {{ book.availableCopies }} экз.</p>
            </li>
        </ul>
        <p v-else-if="searchPerformed">Книги не найдены.</p>

    </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';

export default {
    setup() {
        const searchQuery = ref('');
        const books = ref([]);
        const searchPerformed = ref(false);

        const searchBooks = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/books?search=${searchQuery.value}`);
                books.value = response.data;
                searchPerformed.value = true;
            } catch (error) {
                console.error('Ошибка при поиске книг:', error);
            }
        };

        return {
            searchQuery,
            books,
            searchPerformed,
            searchBooks
        };
    },
};
</script>

<style scoped></style>