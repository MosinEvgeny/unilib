<template>
    <div class="container">
        <h1>Управление книгами</h1>

        <h3>Добавить новую книгу</h3>
        <form @submit.prevent="addBook">
            <div class="form-group">
                <label for="title">Название:</label>
                <input type="text" id="title" v-model="newBook.title" required>
            </div>
            <div class="form-group">
                <label for="author">Автор:</label>
                <input type="text" id="author" v-model="newBook.author" required>
            </div>
            <div class="form-group">
                <label for="isbn">ISBN:</label>
                <input type="text" id="isbn" v-model="newBook.isbn">
            </div>
            <div class="form-group">
                <label for="publisher">Издательство:</label>
                <input type="text" id="publisher" v-model="newBook.publisher">
            </div>
            <div class="form-group">
                <label for="publication_year">Год издания:</label>
                <input type="number" id="publication_year" v-model="newBook.publication_year">
            </div>
            <div class="form-group">
                <label for="total_copies">Количество экземпляров:</label>
                <input type="number" id="total_copies" v-model="newBook.total_copies" min="1" required>
            </div>
            <div class="form-group">
                <label for="category">Категория:</label>
                <input type="text" id="category" v-model="newBook.category">
            </div>
            <div class="form-group">
                <label for="description">Описание:</label>
                <textarea id="description" v-model="newBook.description"></textarea>
            </div>
            <button type="submit" :disabled="isAdding">Добавить</button>
            <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
            <p v-if="successMessage" class="success">{{ successMessage }}</p>
        </form>

        <h3>Редактировать книгу</h3>
        <div class="form-group">
            <label for="searchQuery">Поиск книги:</label>
            <input type="text" id="searchQuery" v-model="searchQuery" placeholder="Название, автор или ISBN">
            <button @click="searchBooks">Поиск</button>
        </div>

        <div v-if="searchResults.length > 0">
            <h3>Результаты поиска:</h3>
            <ul>
                <li v-for="book in searchResults" :key="book.book_id">
                    {{ book.title }} - {{ book.author }} (ISBN: {{ book.isbn }})
                    <button @click="selectBookForEdit(book)">Загрузить книгу</button>
                </li>
            </ul>
        </div>
        <p v-else-if="searchPerformed">Книги не найдены.</p>

        <form v-if="editingBook" @submit.prevent="updateBook">
            <div class="form-group">
                <label for="editTitle">Название:</label>
                <input type="text" id="editTitle" v-model="editingBook.title" required>
            </div>
            <div class="form-group">
                <label for="editAuthor">Автор:</label>
                <input type="text" id="editAuthor" v-model="editingBook.author" required>
            </div>
            <div class="form-group">
                <label for="editIsbn">ISBN:</label>
                <input type="text" id="editIsbn" v-model="editingBook.isbn">
            </div>
            <div class="form-group">
                <label for="editPublisher">Издательство:</label>
                <input type="text" id="editPublisher" v-model="editingBook.publisher">
            </div>
            <div class="form-group">
                <label for="editPublicationYear">Год издания:</label>
                <input type="number" id="editPublicationYear" v-model="editingBook.publication_year">
            </div>
            <div class="form-group">
                <label for="editTotalCopies">Количество экземпляров:</label>
                <input type="number" id="editTotalCopies" v-model="editingBook.total_copies" min="1" required>
            </div>
            <div class="form-group">
                <label for="editCategory">Категория:</label>
                <input type="text" id="editCategory" v-model="editingBook.category">
            </div>
            <div class="form-group">
                <label for="editDescription">Описание:</label>
                <textarea id="editDescription" v-model="editingBook.description"></textarea>
            </div>
            <button type="submit" :disabled="isUpdating">Сохранить изменения</button>
        </form>
        <p v-if="editErrorMessage" class="error">{{ editErrorMessage }}</p>
        <p v-if="editSuccessMessage" class="success">{{ editSuccessMessage }}</p>

        <h3>Списать книгу</h3>
        <div class="form-group">
            <label for="searchQueryRemoval">Поиск книги:</label>
            <input type="text" id="searchQueryRemoval" v-model="searchQueryRemoval"
                placeholder="Название, автор или ISBN">
            <button @click="searchBooksForRemoval">Поиск</button>
        </div>

        <div v-if="searchResultsRemoval.length > 0">
            <h3>Результаты поиска:</h3>
            <ul>
                <li v-for="book in searchResultsRemoval" :key="book.book_id">
                    {{ book.title }} - {{ book.author }} (ISBN: {{ book.isbn }})
                    <button @click="selectBookForRemoval(book)">Выбрать для списания</button>
                </li>
            </ul>
        </div>
        <p v-else-if="searchPerformedRemoval">Книги не найдены.</p>

        <div v-if="selectedBookForRemoval" class="book-info">
            <h3>Выбранная книга:</h3>
            <p>Название: {{ selectedBookForRemoval.title }}</p>
            <p>Автор: {{ selectedBookForRemoval.author }}</p>
            <p>ISBN: {{ selectedBookForRemoval.isbn }}</p>
        </div>
        <button v-if="selectedBookForRemoval" @click="removeBook" :disabled="isRemoving">Списать</button>
        <p v-if="removeErrorMessage" class="error">{{ removeErrorMessage }}</p>
        <p v-if="removeSuccessMessage" class="success">{{ removeSuccessMessage }}</p>

    </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';

export default {
    setup() {
        const newBook = ref({
            title: '',
            author: '',
            isbn: '',
            publisher: '',
            publication_year: null,
            total_copies: 1,
            category: '',
            description: ''
        });

        const editBookId = ref(null);
        const editingBook = ref(null);
        const editErrorMessage = ref('');
        const editSuccessMessage = ref('');
        const isUpdating = ref(false);

        const errorMessage = ref('');
        const successMessage = ref('');
        const isAdding = ref(false);

        const removeBookId = ref(null);
        const removeErrorMessage = ref('');
        const removeSuccessMessage = ref('');
        const isRemoving = ref(false);

        const books = ref([]);

        const searchQuery = ref('');
        const searchResults = ref([]);
        const searchPerformed = ref(false);

        const searchQueryRemoval = ref('');
        const searchResultsRemoval = ref([]);
        const searchPerformedRemoval = ref(false);

        const selectedBookForRemoval = ref(null);

        const fetchAllBooks = async () => {
            try {
                const response = await axios.get('http://localhost:8080/books');
                books.value = response.data;
            } catch (error) {
                console.error('Ошибка при загрузке книг:', error);
            }
        };

        const searchBooks = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/books?search=${searchQuery.value}`);
                searchResults.value = response.data;
                searchPerformed.value = true;
            } catch (error) {
                console.error('Ошибка при поиске книг:', error);
            }
        };

        const searchBooksForRemoval = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/books?search=${searchQueryRemoval.value}`);
                searchResultsRemoval.value = response.data;
                searchPerformedRemoval.value = true;
            } catch (error) {
                console.error('Ошибка при поиске книг:', error);
            }
        };

        const selectBookForEdit = (book) => {
            editingBook.value = { ...book };
            searchResults.value = [];
            searchQuery.value = '';
            searchPerformed.value = false;
        };

        const selectBookForRemoval = (book) => {
            selectedBookForRemoval.value = book;
            searchResultsRemoval.value = [];
            searchQueryRemoval.value = '';
            searchPerformedRemoval.value = false;
        };

        const removeBook = async () => {
            try {
                isRemoving.value = true;
                const response = await axios.delete(`http://localhost:8080/books/${selectedBookForRemoval.value.book_id}`);

                if (response.status === 200) {
                    removeSuccessMessage.value = 'Книга  успешно  списана.';
                    removeErrorMessage.value = '';
                    selectedBookForRemoval.value = null;
                    //  Обновить список книг после списания
                    await fetchAllBooks();
                } else {
                    const errorData = await response.json();
                    removeErrorMessage.value = errorData.error || 'Произошла ошибка при списании книги.';
                }
            } catch (error) {
                console.error('Ошибка  при  списании  книги:', error);
                removeErrorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            } finally {
                isRemoving.value = false;
            }
        };

        const fetchBook = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/books/${editBookId.value}`);
                editingBook.value = response.data;
                editErrorMessage.value = '';
            } catch (error) {
                editingBook.value = null;
                editErrorMessage.value = 'Книга не найдена.';
            }
        };

        const updateBook = async () => {
            try {
                isUpdating.value = true;
                const response = await axios.put(`http://localhost:8080/books/${editingBook.value.book_id}`, editingBook.value);

                if (response.status === 200) {
                    editSuccessMessage.value = 'Изменения успешно сохранены.';
                    editErrorMessage.value = '';
                } else {
                    const errorData = await response.json();
                    editErrorMessage.value = errorData.error || 'Произошла ошибка при сохранении изменений.';
                }
            } catch (error) {
                console.error('Ошибка при редактировании книги:', error);
                editErrorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            } finally {
                isUpdating.value = false;
            }
        };

        const addBook = async () => {
            try {
                isAdding.value = true;
                const response = await axios.post('http://localhost:8080/books', newBook.value);

                if (response.status === 201) {
                    successMessage.value = 'Книга успешно добавлена.';
                    errorMessage.value = '';
                    // Сбросить поля формы
                    newBook.value = {
                        title: '',
                        author: '',
                        isbn: '',
                        publisher: '',
                        publication_year: null,
                        total_copies: 1,
                        category: '',
                        description: ''
                    };

                    // Получить ID новой книги из ответа 
                    const newBookData = response.data;

                    // Добавить новую книгу в массив books
                    books.value.push(newBookData);

                    await fetchAllBooks();
                    
                } else {
                    const errorData = await response.json();
                    errorMessage.value = errorData.error || 'Произошла ошибка при добавлении книги.';
                }
            } catch (error) {
                console.error('Ошибка при добавлении книги:', error);
                errorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            } finally {
                isAdding.value = false;
            }
        };

        onMounted(fetchAllBooks);

        return {
            newBook,
            errorMessage,
            successMessage,
            isAdding,
            addBook,
            editBookId,
            editingBook,
            editErrorMessage,
            editSuccessMessage,
            isUpdating,
            fetchBook,
            updateBook,
            removeBookId,
            removeErrorMessage,
            removeSuccessMessage,
            isRemoving,
            books,
            searchQuery,
            searchResults,
            searchPerformed,
            searchQueryRemoval,
            searchResultsRemoval,
            searchPerformedRemoval,
            selectBookForEdit,
            selectBookForRemoval,
            selectedBookForRemoval,
            removeBook,
            searchBooks,
            searchBooksForRemoval,
        };
    },
};
</script>

<style scoped>
ul {
    list-style: none;
    padding: 0;
}

li {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
    padding: 1rem;
    border: 1px solid #eee;
    border-radius: 4px;
}

button {
    white-space: nowrap;
}
</style>