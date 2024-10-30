<template>
  <div class="container">
    <h1>UniLib</h1>
    <p>Вход в систему</p>
    <form @submit.prevent="login">
      <div class="form-group">
        <label for="username">Логин:</label>
        <input type="text" id="username" v-model="username" required>
      </div>
      <div class="form-group">
        <label for="password">Пароль:</label>
        <input type="password" id="password" v-model="password" required>
      </div>
      <button type="submit">Войти</button>
      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
    </form>
    <p>Нет аккаунта? <router-link to="/register">Зарегистрироваться</router-link></p>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

export default {
  setup() {
    const username = ref('');
    const password = ref('');
    const errorMessage = ref('');
    const router = useRouter();

    const login = async () => {
      try {
        const response = await fetch('http://localhost:8080/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            username: username.value,
            password: password.value
          })
        });

        if (response.ok) {
          const data = await response.json();
          console.log('Ответ сервера:', data);

          localStorage.setItem('role', data.role);
          localStorage.setItem('username', data.username);

          if (data.role === 'admin') {
            router.push('/admin/manage-books');
          } else if (data.role === 'librarian') {
            router.push('/librarian/search-book');
          } else if (data.role === 'reader') {
            localStorage.setItem('reader_id', data.reader_id);
            localStorage.setItem('student_id', data.student_id);
            router.push('/reader');
          } else {
            errorMessage.value = 'Неизвестная  роль  пользователя.';
          }
        } else {
          const errorData = await response.json();
          errorMessage.value = errorData.error || 'Произошла ошибка. Попробуйте позже.';
        }

      } catch (error) {
        console.error('Ошибка при аутентификации:', error);
        errorMessage.value = 'Произошла ошибка. Попробуйте позже.';
      }
    };

    return { username, password, errorMessage, login };
  },
};
</script>

<style scoped>
.container {
  /* ... (другие стили) ... */
  display: flex;
  flex-direction: column;
  align-items: center;
  /* Выравнивание по горизонтали */
  justify-content: center;
  /* Выравнивание по вертикали */
  height: 100vh;
  /* Занимаем всю высоту экрана */
}
</style>