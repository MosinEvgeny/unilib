<template>
    <div class="container">
        <h1>Отчеты</h1>

        <div v-if="showReport">
            <h3>Статистика за последний месяц:</h3>
            <p>Зарегистрированных читателей: {{ reportData.registeredReaders }}</p>
            <p>Книг выдано для чтения: {{ reportData.issuedBooks }}</p>
            <p>Книг возвращено: {{ reportData.returnedBooks }}</p>
            <p>Новых книг: {{ reportData.newBooks }}</p>
            <button @click="generateReportFile" :disabled="isGeneratingReport">
                Сохранить отчет в PDF
            </button>
        </div>

        <button @click="generateReport" :disabled="isFetchingReport">
            Сгенерировать отчет
        </button>

        <p v-if="reportErrorMessage" class="error">{{ reportErrorMessage }}</p>
        <p v-if="reportSuccessMessage" class="success">{{ reportSuccessMessage }}</p>
    </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';

export default {
    setup() {
        const reportErrorMessage = ref('');
        const reportSuccessMessage = ref('');
        const isFetchingReport = ref(false);
        const isGeneratingReport = ref(false);
        const showReport = ref(false);
        const reportData = ref({});

        const generateReport = async () => {
            try {
                isFetchingReport.value = true;
                const response = await axios.get('http://localhost:8080/reports/operations');

                if (response.status === 200) {
                    reportData.value = response.data;
                    showReport.value = true;
                    reportErrorMessage.value = '';
                } else {
                    const errorData = await response.json();
                    reportErrorMessage.value = errorData.error || 'Произошла ошибка при генерации отчета.';
                }
            } catch (error) {
                console.error('Ошибка при генерации отчета:', error);
                reportErrorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            } finally {
                isFetchingReport.value = false;
            }
        };

        const generateReportFile = async () => {
            try {
                isGeneratingReport.value = true;
                const response = await axios.post('http://localhost:8080/reports/operations/generate', {
                    registeredReaders: reportData.value.registeredReaders,
                    issuedBooks: reportData.value.issuedBooks,
                    returnedBooks: reportData.value.returnedBooks,
                    adminName: reportData.value.adminName,
                    librarianName: reportData.value.librarianName,
                    newBooks: reportData.value.newBooks,
                });

                if (response.status === 200) {
                    reportSuccessMessage.value = `Отчет успешно создан: ${response.data.filename}`;
                    reportErrorMessage.value = '';
                } else {
                    const errorData = await response.json();
                    reportErrorMessage.value = errorData.error || 'Произошла ошибка при создании отчета.';
                }
            } catch (error) {
                console.error('Ошибка при генерации отчета:', error);
                reportErrorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            } finally {
                isGeneratingReport.value = false;
            }
        };

        return {
            reportErrorMessage,
            reportSuccessMessage,
            isFetchingReport,
            isGeneratingReport,
            showReport,
            reportData,
            generateReport,
            generateReportFile
        };
    },
};
</script>