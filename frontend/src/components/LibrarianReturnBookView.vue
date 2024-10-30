<template>
    <div class="container">
        <h1>Прием книги</h1>

        <div class="form-group">
            <label for="studentId">Номер студенческого билета:</label>
            <input type="text" id="studentId" v-model="studentId" required>
            <button @click="fetchIssues">Найти выдачи</button>
        </div>

        <div v-if="activeIssues.length > 0">
            <h3>Активные выдачи:</h3>
            <ul class="issue-list">
                <li v-for="issue in activeIssues" :key="issue.issue_id">
                    <p>ID выдачи: {{ issue.issue_id }}</p>
                    <p>Название книги: {{ issue.book_title }}</p>
                    <p>Дата выдачи: {{ formatDate(issue.issue_date) }}</p>
                    <p>Срок возврата: {{ formatDate(issue.due_date) }}</p>
                    <button v-if="!issue.return_date" @click="returnBook(issue.issue_id)"
                        :disabled="isReturning">Принять книгу</button>
                </li>
            </ul>
        </div>
        <p v-else-if="searchPerformed">Активных выдач не найдено.</p>

        <div v-if="returnedIssues.length > 0">
            <h3>История выдач:</h3>
            <ul class="issue-list">
                <li v-for="issue in returnedIssues" :key="issue.issue_id">
                    <p>ID выдачи: {{ issue.issue_id }}</p>
                    <p>Название книги: {{ issue.book_title }}</p>
                    <p>Дата выдачи: {{ formatDate(issue.issue_date) }}</p>
                    <p>Срок возврата: {{ formatDate(issue.due_date) }}</p>
                </li>
            </ul>
        </div>
        <p v-else-if="searchPerformed">Выдачи не найдены.</p>

        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
        <p v-if="successMessage" class="success">{{ successMessage }}</p>
    </div>
</template>

<script>
import { ref, computed } from 'vue';
import axios from 'axios';

export default {
    setup() {
        const studentId = ref('');
        const issues = ref([]);
        const errorMessage = ref('');
        const successMessage = ref('');
        const isReturning = ref(false);
        const searchPerformed = ref(false);
        const activeIssues = computed(() => issues.value.filter(issue => !issue.return_date));
        const returnedIssues = computed(() => issues.value.filter(issue => issue.return_date));

        const fetchIssues = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/reader/${studentId.value}/issues`);
                issues.value = response.data;
                searchPerformed.value = true;
                errorMessage.value = '';
            } catch (error) {
                issues.value = [];
                errorMessage.value = 'Ошибка при поиске выдач.';
                console.error('Ошибка при поиске выдач:', error);
            }
        };

        const returnBook = async (issueId) => {
            try {
                isReturning.value = true;
                const response = await axios.put(`http://localhost:8080/issue/${issueId}/return`);

                if (response.status === 200) {
                    successMessage.value = 'Книга успешно принята.';
                    errorMessage.value = '';
                    //  Обновляем  список  выдач  после  приема  книги 
                    await fetchIssues();
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
            studentId,
            issues,
            errorMessage,
            successMessage,
            isReturning,
            searchPerformed,
            fetchIssues,
            returnBook,
            formatDate,
            activeIssues,
            returnedIssues,
        };
    },
};
</script>

<style scoped></style>