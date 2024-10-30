<template>
    <div class="container">
        <h2>История отчетов</h2>
        <p v-if="deleteSuccessMessage" class="success">{{ deleteSuccessMessage }}</p>
        <ul v-if="reports.length > 0">
            <li v-for="report in reports" :key="report.ReportID">
                {{ report.formatted_date }} - {{ report.ReportContent }}
                <button @click="deleteReport(report.ReportID)">Удалить запись</button>
            </li>
        </ul>
        <p v-else>Отчеты не найдены.</p>
    </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';

export default {
    setup() {
        const reports = ref([]);
        const deleteSuccessMessage = ref('');

        const fetchReports = async () => {
            try {
                const response = await axios.get('http://localhost:8080/reports');
                reports.value = response.data;
            } catch (error) {
                console.error('Ошибка при загрузке отчетов:', error);
            }
        };

        const deleteReport = async (reportId) => {
            try {
                const response = await axios.delete(`http://localhost:8080/reports/${reportId}`);

                if (response.status === 200) {
                    //  Удаляем  отчет  из  массива  reports 
                    reports.value = reports.value.filter(report => report.ReportID !== reportId);
                    //  Отображаем  сообщение  об  успехе 
                    deleteSuccessMessage.value = 'Отчет  успешно  удален!';

                    setTimeout(() => {
                        deleteSuccessMessage.value = '';
                    }, 3000);
                } else {
                    console.error('Ошибка при удалении отчета:', response.data.error);
                }
            } catch (error) {
                console.error('Ошибка при удалении отчета:', error);
            }
        };

        const formatDate = (dateString) => {
            const date = new Date(dateString);
            const options = { year: 'numeric', month: 'long', day: 'numeric' };
            return date.toLocaleDateString('ru-RU', options);
        };

        onMounted(fetchReports);

        return {
            reports,
            formatDate,
            deleteSuccessMessage,
            deleteReport,
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