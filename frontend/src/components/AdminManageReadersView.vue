<template>
    <div class="container">
        <h1>Управление читателями</h1>

        <div class="search-filters">
            <input type="text" v-model="searchQuery" placeholder="Поиск по ФИО или номеру студенческого...">
            <button @click="searchReaders">Поиск</button>
        </div>

        <h3>Список читателей</h3>
        <li v-for="reader in readers" :key="reader.reader_id">
            {{ reader.full_name }} - {{ reader.student_id }}
            <button @click="editReader(reader)">Редактировать</button>
            <button @click="deleteReader(reader.reader_id)">Удалить</button>
        </li>

        <form v-if="editingReader" @submit.prevent="updateReader">
            <h3>Редактировать читателя</h3>
            <div class="form-group">
                <label for="editFullName">ФИО:</label>
                <input type="text" id="editFullName" v-model="editingReader.full_name" required>
            </div>
            <div class="form-group">
                <label for="editfaculty">Факультет:</label>
                <select id="editfaculty" v-model="editingReader.faculty">
                    <option value="">Выберите факультет</option>
                    <option value="Математический">Математический</option>
                    <option value="Филологический">Филологический</option>
                    <option value="Физической культуры">Физической культуры</option>
                    <option value="Педагогики и методики начального образования">Педагогики и методики начального
                        образования
                    </option>
                    <option value="Педагогики и психологии детства">Педагогики и психологии детства</option>
                    <option value="Естественнонаучный">Естественнонаучный</option>
                    <option value="Иностранных языков">Иностранных языков</option>
                    <option value="Физический">Физический</option>
                    <option value="Исторический">Исторический</option>
                    <option value="Музыки">Музыки</option>
                    <option value="Психологии">Психологии</option>
                    <option value="Информатики и экономики">Информатики и экономики</option>
                    <option value="Правового и социально-педагогического образования">Правового и
                        социально-педагогического
                        образования</option>
                    <option value="Кафедра педагогики и психологии">Кафедра педагогики и психологии</option>
                    <option value="Отдел подготовки научно-педагогических кадров">Отдел подготовки научно-педагогических
                        кадров
                    </option>
                    <option value="Международное образование">Международное образование</option>
                    <option value="Открытый университет">Открытый университет</option>
                </select>
            </div>
            <div class="form-group">
                <label for="editcourse">Курс:</label>
                <input type="number" id="editcourse" v-model="editingReader.course" min="1" max="6">
            </div>
            <div class="form-group">
                <label for="editStudentId">Номер студенческого билета:</label>
                <input type="text" id="editStudentId" v-model="editingReader.student_id">
            </div>
            <div class="form-group">
                <label for="editphone_number">Номер телефона:</label>
                <input type="text" id="editphone_number" v-model="editingReader.phone_number">
            </div>
            <div class="form-group">
                <label for="editusername">Логин:</label>
                <input type="text" id="editusername" v-model="editingReader.username" required>
            </div>
            <div class="form-group">
                <label for="editpassword">Пароль:</label>
                <input type="text" id="editpassword" v-model="editingReader.password" required>
            </div>
            <button type="submit" :disabled="isUpdatingReader">Сохранить</button>
        </form>
        <p v-if="readerEditErrorMessage" class="error">{{ readerEditErrorMessage }}</p>
        <p v-if="readerEditSuccessMessage" class="success">{{ readerEditSuccessMessage }}</p>

        <p v-if="readerDeleteErrorMessage" class="error">{{ readerDeleteErrorMessage }}</p>
        <p v-if="readerDeleteSuccessMessage" class="success">{{ readerDeleteSuccessMessage }}</p>

    </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';

export default {
    setup() {
        const readers = ref([]);
        const editingReader = ref(null);
        const isUpdatingReader = ref(false);
        const readerEditErrorMessage = ref('');
        const readerEditSuccessMessage = ref('');
        const readerDeleteErrorMessage = ref('');
        const readerDeleteSuccessMessage = ref('');
        const searchQuery = ref('');

        const searchReaders = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/readers?search=${searchQuery.value}`);
                readers.value = response.data;
            } catch (error) {
                console.error('Ошибка  при  поиске  читателей:', error);
            }
        };

        const fetchReaders = async () => {
            try {
                const response = await axios.get('http://localhost:8080/readers');
                readers.value = response.data;
            } catch (error) {
                console.error('Ошибка при загрузке читателей:', error);
            }
        };

        const editReader = (reader) => {
            editingReader.value = { ...reader };
        };

        const updateReader = async () => {
            try {
                isUpdatingReader.value = true;
                const response = await axios.put(`http://localhost:8080/readers/${editingReader.value.reader_id}`, editingReader.value);

                if (response.status === 200) {
                    readerEditSuccessMessage.value = 'Данные читателя успешно обновлены.';
                    readerEditErrorMessage.value = '';
                    // Обнови данные в массиве readers
                    const readerIndex = readers.value.findIndex(r => r.reader_id === editingReader.value.reader_id);
                    if (readerIndex !== -1) {
                        readers.value.splice(readerIndex, 1, editingReader.value);
                    }
                    editingReader.value = null;
                } else {
                    const errorData = await response.json();
                    readerEditErrorMessage.value = errorData.error || 'Произошла ошибка при обновлении данных.';
                }
            } catch (error) {
                console.error('Ошибка при обновлении данных читателя:', error);
                readerEditErrorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            } finally {
                isUpdatingReader.value = false;
            }
        };

        const deleteReader = async (readerId) => {
            try {
                const response = await axios.delete(`http://localhost:8080/readers/${readerId}`);

                if (response.status === 200) {
                    readerDeleteSuccessMessage.value = 'Читатель успешно удален.';
                    readerDeleteErrorMessage.value = '';
                    // Удали читателя из массива readers
                    readers.value = readers.value.filter(reader => reader.reader_id !== readerId);
                } else {
                    const errorData = await response.json();
                    readerDeleteErrorMessage.value = errorData.error || 'Произошла ошибка при удалении читателя.';
                }
            } catch (error) {
                console.error('Ошибка при удалении читателя:', error);
                readerDeleteErrorMessage.value = 'Произошла ошибка. Попробуйте позже.';
            }
        };

        onMounted(fetchReaders);

        return {
            readers,
            editingReader,
            isUpdatingReader,
            readerEditErrorMessage,
            readerEditSuccessMessage,
            readerDeleteErrorMessage,
            readerDeleteSuccessMessage,
            fetchReaders,
            editReader,
            updateReader,
            deleteReader,
            searchQuery,
            searchReaders,
        };
    },
};
</script>

<style scoped></style>