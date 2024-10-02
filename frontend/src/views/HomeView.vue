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
    <router-link to="/register">Зарегистрироваться</router-link>
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

          if (data.role === 'admin') {
            localStorage.setItem('role', data.role);
            router.push('/admin');
          } else if (data.role === 'librarian') {
            localStorage.setItem('role', data.role);
            router.push('/librarian');
          } else if (data.role === 'reader') {
            localStorage.setItem('role', data.role);
            localStorage.setItem('reader_id', data.reader_id);
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

<style scoped></style>