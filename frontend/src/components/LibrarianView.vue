<template>
    <div class="container">
        <h1>UniLib - Панель библиотекаря</h1>

        <nav>
            <router-link to="/librarian/register-reader">Регистрация читателя</router-link>
            <router-link to="/librarian/issue-book">Выдача книги</router-link>
            <router-link to="/librarian/return-book">Прием книги</router-link>
            <router-link to="/librarian/search-book">Поиск книг</router-link>
        </nav>

        <router-view />

        <button @click="logout">Выход</button>
    </div>
</template>

<script>
import { useRouter } from 'vue-router';

export default {
    setup() {
        const router = useRouter();

        const logout = () => {
            // Логика выхода (очистка данных авторизации, перенаправление на страницу входа)
            router.push('/');
        };

        return { logout };
    },
};
</script>

<style scoped>
/* Стили для панели библиотекаря */
</style>