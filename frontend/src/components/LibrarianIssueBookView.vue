<template>
    <div class="container">
        <h2>Выдача книги</h2>

        <div class="form-group">
            <label for="studentId">Номер студенческого билета:</label>
            <input type="text" id="studentId" v-model="studentId" required>
            <button @click="fetchReader">Найти читателя</button>
        </div>

        <div v-if="reader" class="reader-info">
            <h3>Информация о читателе:</h3>
            <p>ФИО: {{ reader.full_name }}</p>
            <p>Факультет: {{ reader.faculty }}</p>
        </div>

        <div class="form-group">
            <label for="bookTitle">Название книги:</label>
            <input type="text" id="bookTitle" v-model="bookTitle" required>
            <button @click="searchBooks">Найти книги</button>
        </div>

        <div v-if="searchResults.length > 0">
            <h3>Результаты поиска:</h3>
            <ul>
                <li v-for="book in searchResults" :key="book.book_id">
                    {{ book.title }} - {{ book.author }}
                    <button @click="selectBook(book)">Выбрать</button>
                </li>
            </ul>
        </div>

        <div v-if="selectedBook" class="book-info">
            <h3>Информация о книге:</h3>
            <p>Название: {{ selectedBook.title }}</p>
            <p>Автор: {{ selectedBook.author }}</p>
            <p v-if="selectedBook.availableCopies > 0">Доступно: {{ selectedBook.availableCopies }} экз.</p>
            <p v-else>Нет в наличии</p>
        </div>

        <button v-if="reader && selectedBook && selectedBook.availableCopies > 0" @click="issueBook"
            :disabled="isIssuing">
            Выдать книгу
        </button>
        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
        <p v-if="successMessage" class="success">{{ successMessage }}</p>

        <router-link to="/librarian">Назад</router-link>
    </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';

export default {
    setup() {
        const studentId = ref('');
        const bookTitle = ref('');
        const searchResults = ref([]);
        const readerId = ref('');
        const bookId = ref('');
        const reader = ref(null);
        const selectedBook = ref(null);
        const errorMessage = ref('');
        const successMessage = ref('');
        const isIssuing = ref(false);

        const fetchReader = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/readers/by-student-id/${studentId.value}`);
                reader.value = response.data;
                errorMessage.value = '';
            } catch (error) {
                reader.value = null;
                errorMessage.value = 'Читатель не найден.';
            }
        };

        const searchBooks = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/books?search=${bookTitle.value}`);
                searchResults.value = response.data;
                errorMessage.value = '';
            } catch (error) {
                // ... (обработка  ошибки) ... 
            }
        };

        const selectBook = async (book) => { //  Добавили  async
            selectedBook.value = book;
            searchResults.value = [];

            //  Получаем  количество  доступных  экземпляров  для  выбранной  книги 
            try {
                const response = await axios.get(`http://localhost:8080/books/${selectedBook.value.book_id}/available`);
                selectedBook.value.availableCopies = response.data.available_copies;
            } catch (error) {
                console.error(`Ошибка при получении доступных экземпляров для книги ${selectedBook.value.book_id}:`, error);
            }
        };

        const fetchBook = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/readers/by-student-id/${studentId.value}`);
                reader.value = response.data;
                errorMessage.value = '';
            } catch (error) {
                book.value = null;
                errorMessage.value = 'Книга не найдена.';
            }
        };

        const issueBook = async () => {
            try {
                isIssuing.value = true;
                const response = await axios.post('http://localhost:8080/issue', {
                    reader_id: reader.value.reader_id,
                    book_title: selectedBook.value.title
                });

                if (response.status === 201) {
                    successMessage.value = 'Книга успешно выдана.';
                    errorMessage.value = '';
                    readerId.value = '';
                    bookId.value = '';
                    reader.value = null;
                    selectedBook.value = null;
                } else {
                    const errorData = await response.json();
                    errorMessage.value = errorData.error || 'Произошла ошибка при выдаче.';
                }
            } catch (error) {
                console.error('Ошибка при выдаче книги:', error);
                errorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            } finally {
                isIssuing.value = false;
            }
        };

        return {
            studentId,
            bookTitle,
            searchResults,
            reader,
            selectedBook,
            errorMessage,
            successMessage,
            isIssuing,
            fetchReader,
            searchBooks,
            selectBook,
            issueBook

        };
    },
};
</script>

<style scoped></style>