<template>
    <div class="container">
        <h2>Прием книги</h2>

        <div class="form-group">
            <label for="issueId">ID выдачи:</label>
            <input type="text" id="issueId" v-model="issueId" required>
            <button @click="fetchIssue">Найти выдачу</button>
        </div>

        <div v-if="issue" class="issue-info">
            <h3>Информация о выдаче:</h3>
            <p>ID читателя: {{ issue.reader_id }}</p>
            <p>ID книги: {{ issue.book_id }}</p>
            <p>Дата выдачи: {{ formatDate(issue.issue_date) }}</p>
            <p>Срок возврата: {{ formatDate(issue.due_date) }}</p>
            <p v-if="issue.return_date">Книга уже возвращена: {{ formatDate(issue.return_date) }}</p>
        </div>

        <button v-if="issue && !issue.return_date" @click="returnBook" :disabled="isReturning">
            Принять книгу
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
        const issueId = ref('');
        const issue = ref(null);
        const errorMessage = ref('');
        const successMessage = ref('');
        const isReturning = ref(false);

        const fetchIssue = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/issue/${issueId.value}`);
                issue.value = response.data;
                errorMessage.value = '';
            } catch (error) {
                issue.value = null;
                errorMessage.value = 'Выдача не найдена.';
            }
        };

        const returnBook = async () => {
            try {
                isReturning.value = true;
                const response = await axios.put(`http://localhost:8080/issue/${issueId.value}/return`);

                if (response.status === 200) {
                    successMessage.value = 'Книга успешно принята.';
                    errorMessage.value = '';
                    issue.value = null; //  Сбросить информацию о выдаче 
                    issueId.value = '';
                } else {
                    const errorData = await response.json();
                    errorMessage.value = errorData.error || 'Произошла ошибка при приеме книги.';
                }
            } catch (error) {
                console.error('Ошибка при приеме книги:', error);
                errorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            } finally {
                isReturning.value = false;
            }
        };

        const formatDate = (dateString) => {
            const date = new Date(dateString);
            return date.toLocaleDateString();
        };

        return {
            issueId,
            issue,
            errorMessage,
            successMessage,
            isReturning,
            fetchIssue,
            returnBook,
            formatDate
        };
    },
};
</script>

<style scoped>

</style>