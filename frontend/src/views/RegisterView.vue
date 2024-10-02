<template>
  <div class="container">
    <h1>UniLib</h1>
    <p>Регистрация</p>
    <form @submit.prevent="register">
      <div class="form-group">
        <label for="fullName">ФИО:</label>
        <input type="text" id="fullName" v-model="fullName" required>
      </div>
      <div class="form-group">
        <label for="faculty">Факультет:</label>
        <select id="faculty" v-model="faculty">
          <option value="">Выберите факультет</option>
          <option value="Математический">Математический</option>
          <option value="Филологический">Филологический</option>
          <option value="Физической культуры">Физической культуры</option>
          <option value="Педагогики и методики начального образования">Педагогики и методики начального образования
          </option>
          <option value="Педагогики и психологии детства">Педагогики и психологии детства</option>
          <option value="Естественнонаучный">Естественнонаучный</option>
          <option value="Иностранных языков">Иностранных языков</option>
          <option value="Физический">Физический</option>
          <option value="Исторический">Исторический</option>
          <option value="Музыки">Музыки</option>
          <option value="Психологии">Психологии</option>
          <option value="Информатики и экономики">Информатики и экономики</option>
          <option value="Правового и социально-педагогического образования">Правового и социально-педагогического
            образования</option>
          <option value="Кафедра педагогики и психологии">Кафедра педагогики и психологии</option>
          <option value="Отдел подготовки научно-педагогических кадров">Отдел подготовки научно-педагогических кадров
          </option>
          <option value="Международное образование">Международное образование</option>
          <option value="Открытый университет">Открытый университет</option>
        </select>
      </div>
      <div class="form-group">
        <label for="course">Курс:</label>
        <input type="number" id="course" v-model="course" min="1" max="6">
      </div>
      <div class="form-group">
        <label for="studentId">Номер студенческого билета:</label>
        <input type="text" id="studentId" v-model="studentId" required>
      </div>
      <div class="form-group">
        <label for="phone_number">Номер телефона:</label>
        <input type="text" id="phone_number" v-model="phone_number">
      </div>
      <div class="form-group">
        <label for="username">Логин:</label>
        <input type="text" id="username" v-model="username" required>
      </div>
      <div class="form-group">
        <label for="password">Пароль:</label>
        <input type="password" id="password" v-model="password" required>
      </div>
      <button type="submit">Зарегистрироваться</button>
      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
    </form>
    <router-link to="/">Уже есть аккаунт? Войти</router-link>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

export default {
  setup() {
    const fullName = ref('');
    const faculty = ref('');
    const course = ref(1);
    const studentId = ref('');
    const phone_number = ref('');
    const username = ref('');
    const password = ref('');
    const errorMessage = ref('');
    const router = useRouter();

    const register = async () => {
      try {
        const response = await fetch('http://localhost:8080/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            full_name: fullName.value,
            faculty: faculty.value,
            course: course.value,
            student_id: studentId.value,
            phone_number: phone_number.value,
            username: username.value,
            password: password.value
          })
        });

        if (response.ok) {
          // Регистрация успешна, перенаправить на страницу входа
          router.push('/');
        } else {
          const errorData = await response.json();
          errorMessage.value = errorData.error || 'Произошла ошибка при регистрации.';
        }
      } catch (error) {
        console.error('Ошибка при регистрации:', error);
        errorMessage.value = 'Произошла ошибка. Попробуйте позже.';
      }
    };

    return {
      fullName,
      faculty,
      course,
      studentId,
      phone_number,
      username,
      password,
      errorMessage,
      register
    };
  },
};
</script>

<style scoped></style>